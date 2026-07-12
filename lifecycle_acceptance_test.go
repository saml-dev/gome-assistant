package gomeassistant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestNewAppDoesNotMakeNetworkRequests(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		requests.Add(1)
	}))
	defer server.Close()

	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         server.URL,
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	require.Zero(t, requests.Load())
	app.Cleanup()
	require.Zero(t, requests.Load())
}

func TestCleanupWaitsForDeliveredCallbackAndRejectsLaterCallback(t *testing.T) {
	callbackStarted := make(chan struct{})
	releaseCallback := make(chan struct{})
	var callbackCount atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/states/zone.home" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"entity_id":"zone.home","state":"home","attributes":{"latitude":1,"longitude":2}}`))
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
		defer conn.Close()
		_ = conn.WriteJSON(map[string]string{"type": "auth_required"})
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		_ = conn.WriteJSON(map[string]string{"type": "auth_ok"})
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var request struct {
				ID        int64  `json:"id"`
				EventType string `json:"event_type"`
			}
			if json.Unmarshal(message, &request) != nil || request.ID == 0 {
				continue
			}
			_ = conn.WriteJSON(map[string]any{"id": request.ID, "type": "result", "success": true, "result": nil})
			if request.EventType == "state_changed" {
				_ = conn.WriteJSON(map[string]any{
					"id":   request.ID,
					"type": "event",
					"event": map[string]any{
						"event_type": "state_changed",
						"data": map[string]any{
							"entity_id": "light.test",
							"old_state": map[string]any{"state": "off"},
							"new_state": map[string]any{"state": "on"},
						},
					},
				})
			}
		}
	}))
	defer server.Close()

	app := fakeApp(t, server.URL)
	app.RegisterEntityListeners(NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {
		if callbackCount.Add(1) == 1 {
			close(callbackStarted)
			<-releaseCallback
		}
	}).Build())

	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(context.Background()) }()
	select {
	case <-callbackStarted:
	case <-time.After(time.Second):
		t.Fatal("state-change callback was not delivered")
	}

	app.Cleanup()
	require.ErrorIs(t, app.ctx.Err(), context.Canceled)

	// This is the same delivered entity payload after cancellation; no second
	// callback may start while Start waits for the first callback to finish.
	app.callEntityListeners([]byte(`{"event":{"data":{"entity_id":"light.test","old_state":{"state":"off"},"new_state":{"state":"on"}}}}`))
	require.Equal(t, int32(1), callbackCount.Load())
	select {
	case <-runDone:
		t.Fatal("Start returned while the delivered callback was blocked")
	case <-time.After(20 * time.Millisecond):
	}

	close(releaseCallback)
	select {
	case err := <-runDone:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Start did not return after Cleanup")
	}
	require.Equal(t, int32(1), callbackCount.Load())
}
