package gomeassistant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"saml.dev/gome-assistant/internal/services"
)

func TestCallCanceledBeforeSend(t *testing.T) {
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = app.Call(ctx, services.BaseServiceRequest{}, nil)
	require.ErrorIs(t, err, context.Canceled)
}

func TestCallReturnsWhenCleanupClosesConnection(t *testing.T) {
	callSeen := make(chan struct{})
	ready := make(chan struct{})
	var callOnce sync.Once
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/states/zone.home" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"entity_id":"zone.home","state":"zoning","attributes":{"latitude":1,"longitude":2}}`))
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
				ID   int64  `json:"id"`
				Type string `json:"type"`
			}
			if json.Unmarshal(message, &request) == nil && request.Type == "call_service" {
				callOnce.Do(func() { close(callSeen) })
			}
			if request.Type == "subscribe_events" {
				select {
				case <-ready:
				default:
					close(ready)
				}
			}
			if request.ID != 0 && request.Type != "call_service" {
				_ = conn.WriteJSON(map[string]any{"id": request.ID, "type": "result", "success": true, "result": nil})
			}
		}
	}))
	t.Cleanup(server.Close)

	app := fakeApp(t, server.URL)
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(context.Background()) }()
	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Fatal("fake server did not authenticate")
	}
	callDone := make(chan error, 1)
	go func() { callDone <- app.Call(context.Background(), services.BaseServiceRequest{}, nil) }()

	select {
	case <-callSeen:
	case <-time.After(time.Second):
		t.Fatal("fake server did not receive call")
	}

	// The observed call is the pending Call's request. Cleaning up the app must
	// wake it through the captured connection's terminal channel.
	app.Cleanup()
	select {
	case err := <-callDone:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("Call did not return after app cleanup")
	}
	select {
	case <-runDone:
	case <-time.After(time.Second):
		t.Fatal("Start did not return after app cleanup")
	}
}

func TestCallReturnsWhenStartupFailureClosesConnectionBeforeStart(t *testing.T) {
	startupLoadStarted := make(chan struct{})
	failStartup := make(chan struct{})
	callSeen := make(chan struct{})
	var startupLoadOnce sync.Once
	var failStartupOnce sync.Once
	var callOnce sync.Once

	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/states/zone.home" {
			startupLoadOnce.Do(func() { close(startupLoadStarted) })
			select {
			case <-failStartup:
				http.Error(w, "startup failed", http.StatusInternalServerError)
			case <-r.Context().Done():
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
		defer conn.Close()
		if err := conn.WriteJSON(map[string]string{"type": "auth_required"}); err != nil {
			return
		}
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		if err := conn.WriteJSON(map[string]string{"type": "auth_ok"}); err != nil {
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var request struct {
				Type string `json:"type"`
			}
			if json.Unmarshal(message, &request) == nil && request.Type == "call_service" {
				callOnce.Do(func() { close(callSeen) })
			}
		}
	}))
	t.Cleanup(func() {
		failStartupOnce.Do(func() { close(failStartup) })
		server.Close()
	})

	app := fakeApp(t, server.URL)
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(context.Background()) }()

	select {
	case <-startupLoadStarted:
	case <-time.After(time.Second):
		t.Fatal("Start did not begin loading startup state")
	}

	callDone := make(chan error, 1)
	go func() { callDone <- app.Call(context.Background(), services.BaseServiceRequest{}, nil) }()
	select {
	case <-callSeen:
	case <-time.After(time.Second):
		t.Fatal("Call was not admitted during startup")
	}

	failStartupOnce.Do(func() { close(failStartup) })
	select {
	case err := <-callDone:
		require.ErrorIs(t, err, ErrConnectionClosed)
	case <-time.After(time.Second):
		t.Fatal("Call did not return after startup closed the connection")
	}
	select {
	case err := <-runDone:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("Start did not return after startup failure")
	}
}
