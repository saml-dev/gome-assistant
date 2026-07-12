package gomeassistant

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func newFakeHA(t *testing.T, stateStatus int) (*httptest.Server, <-chan struct{}, <-chan struct{}) {
	t.Helper()
	ready := make(chan struct{})
	disconnected := make(chan struct{})
	var readyOnce sync.Once
	var disconnectedOnce sync.Once
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/states/") {
			w.WriteHeader(stateStatus)
			if stateStatus == http.StatusOK {
				entityID := strings.TrimPrefix(r.URL.Path, "/api/states/")
				attributes := map[string]any{}
				if entityID == "zone.home" {
					attributes["latitude"] = 1
					attributes["longitude"] = 2
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"entity_id":  entityID,
					"state":      "zoning",
					"attributes": attributes,
				})
			}
			return
		}
		if r.URL.Path != "/api/websocket" {
			http.NotFound(w, r)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() {
			_ = conn.Close()
			disconnectedOnce.Do(func() { close(disconnected) })
		}()
		_ = conn.WriteJSON(map[string]string{"type": "auth_required"})
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		_ = conn.WriteJSON(map[string]string{"type": "auth_ok"})
		readyOnce.Do(func() { close(ready) })
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var request struct {
				ID int64 `json:"id"`
			}
			if json.Unmarshal(message, &request) == nil && request.ID != 0 {
				_ = conn.WriteJSON(map[string]any{"id": request.ID, "type": "result", "success": true, "result": nil})
			}
		}
	}))
	t.Cleanup(server.Close)
	return server, ready, disconnected
}

func fakeApp(t *testing.T, serverURL string) *App {
	t.Helper()
	app, err := NewApp(context.Background(), NewAppRequest{URL: serverURL, HAAuthToken: "token"})
	require.NoError(t, err)
	return app
}

func TestStartAfterPriorCleanupReturnsErrAppClosed(t *testing.T) {
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	app.Cleanup()
	require.ErrorIs(t, app.Start(context.Background()), ErrAppClosed)
}

func TestRunOnStartupRunsOncePerSequentialStartSessions(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)

	calls := make(chan EntityData, 3)
	var callbackCount atomic.Int32
	app.RegisterEntityListeners(NewEntityListener().
		EntityIDs("light.first", "light.second").
		Call(func(_ *Service, _ State, data EntityData) {
			callbackCount.Add(1)
			calls <- data
		}).
		RunOnStartup().
		Build())

	for session := 1; session <= 2; session++ {
		ctx, cancel := context.WithCancel(context.Background())
		runDone := make(chan error, 1)
		go func() { runDone <- app.Start(ctx) }()

		select {
		case data := <-calls:
			require.Contains(t, []string{"light.first", "light.second"}, data.TriggerEntityID)
		case <-time.After(time.Second):
			t.Fatalf("RunOnStartup callback did not run in session %d", session)
		}

		select {
		case data := <-calls:
			t.Fatalf("RunOnStartup callback ran twice in session %d: %s", session, data.TriggerEntityID)
		case <-time.After(50 * time.Millisecond):
		}

		cancel()
		select {
		case err := <-runDone:
			require.ErrorIs(t, err, context.Canceled)
		case <-time.After(time.Second):
			t.Fatalf("Start did not return after session %d cancellation", session)
		}
		require.Equal(t, int32(session), callbackCount.Load())
	}
}

func TestSolarScheduleInitializationIsDeferredUntilHomeZoneLoad(t *testing.T) {
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)

	schedule := NewDailySchedule().Call(func(*Service, State) {}).Sunrise().Build()
	app.RegisterSchedules(schedule)
	registered, ok := app.scheduledActions[0].(*DailySchedule)
	require.True(t, ok)
	require.True(t, registered.nextRunTime.IsZero())

	// This is the state established by loadHomeZone; calling the focused
	// initializer here keeps the regression test network-free.
	app.state.latitude = 40.7128
	app.state.longitude = -74.0060
	app.initializeSolarSchedules()
	initialized, ok := app.scheduledActions[0].(*DailySchedule)
	require.True(t, ok)
	require.False(t, initialized.nextRunTime.IsZero())
}

func TestStartCancellationUnblocks(t *testing.T) {
	server, ready, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- app.Start(ctx) }()
	<-ready
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
}

func TestCleanupIsIdempotent(t *testing.T) {
	server, ready, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	done := make(chan error, 1)
	go func() { done <- app.Start(context.Background()) }()
	<-ready

	app.Cleanup()
	require.ErrorIs(t, <-done, context.Canceled)
	app.Cleanup()
}

func TestCleanupCancelsImmediatelyStartingApp(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)

	const iterations = 200
	for i := 0; i < iterations; i++ {
		app := fakeApp(t, server.URL)
		runDone := make(chan error, 1)
		go func() { runDone <- app.Start(context.Background()) }()

		app.Cleanup()
		select {
		case err := <-runDone:
			require.True(t, errors.Is(err, ErrAppClosed) || errors.Is(err, context.Canceled), "iteration %d: %v", i, err)
		case <-time.After(time.Second):
			t.Fatalf("iteration %d: Start did not return after Cleanup", i)
		}
	}
}

func TestStartFailureClosesSocket(t *testing.T) {
	server, ready, disconnected := newFakeHA(t, http.StatusNotFound)
	app := fakeApp(t, server.URL)
	done := make(chan error, 1)
	go func() { done <- app.Start(context.Background()) }()
	<-ready
	require.Error(t, <-done)
	select {
	case <-disconnected:
	case <-time.After(time.Second):
		t.Fatal("Start did not close the websocket after startup failure")
	}
	app.Cleanup()
}
