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
)

func TestStartCancellationCancelsStalledHomeZoneRequest(t *testing.T) {
	for _, stop := range []struct {
		name string
		call func(context.CancelFunc, *App) error
	}{
		{
			name: "parent context",
			call: func(cancel context.CancelFunc, _ *App) error {
				cancel()
				return nil
			},
		},
		{
			name: "app cleanup",
			call: func(_ context.CancelFunc, app *App) error {
				app.Cleanup()
				return nil
			},
		},
	} {
		t.Run(stop.name, func(t *testing.T) {
			requestStarted := make(chan struct{})
			requestCanceled := make(chan struct{})
			var startedOnce sync.Once
			var canceledOnce sync.Once
			upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/states/zone.home":
					startedOnce.Do(func() { close(requestStarted) })
					<-r.Context().Done()
					canceledOnce.Do(func() { close(requestCanceled) })
				case "/api/websocket":
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
					_, _, _ = conn.ReadMessage()
				default:
					http.NotFound(w, r)
				}
			}))
			t.Cleanup(server.Close)

			app, err := NewApp(context.Background(), NewAppRequest{
				URL:         server.URL,
				HAAuthToken: "token",
			})
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			runDone := make(chan error, 1)
			go func() { runDone <- app.Start(ctx) }()

			select {
			case <-requestStarted:
			case <-time.After(time.Second):
				t.Fatal("Start did not request the home zone")
			}
			require.NoError(t, stop.call(cancel, app))

			select {
			case err := <-runDone:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(time.Second):
				t.Fatal("Start did not return after cancellation")
			}
			select {
			case <-requestCanceled:
			case <-time.After(time.Second):
				t.Fatal("home-zone request was not canceled")
			}
		})
	}
}

func TestStartCancellationCancelsStalledRunOnStartupStateRequest(t *testing.T) {
	for _, stop := range []struct {
		name string
		call func(context.CancelFunc, *App) error
	}{
		{
			name: "parent context",
			call: func(cancel context.CancelFunc, _ *App) error {
				cancel()
				return nil
			},
		},
		{
			name: "app cleanup",
			call: func(_ context.CancelFunc, app *App) error {
				app.Cleanup()
				return nil
			},
		},
	} {
		t.Run(stop.name, func(t *testing.T) {
			requestStarted := make(chan struct{})
			requestCanceled := make(chan struct{})
			var startedOnce sync.Once
			var canceledOnce sync.Once
			upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/states/zone.home":
					_ = json.NewEncoder(w).Encode(map[string]any{
						"entity_id":  "zone.home",
						"state":      "zoning",
						"attributes": map[string]any{"latitude": 1, "longitude": 2},
					})
				case "/api/states/light.startup":
					startedOnce.Do(func() { close(requestStarted) })
					<-r.Context().Done()
					canceledOnce.Do(func() { close(requestCanceled) })
				case "/api/websocket":
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
							ID int64 `json:"id"`
						}
						if json.Unmarshal(message, &request) == nil && request.ID != 0 {
							_ = conn.WriteJSON(map[string]any{"id": request.ID, "type": "result", "success": true, "result": nil})
						}
					}
				default:
					http.NotFound(w, r)
				}
			}))
			t.Cleanup(server.Close)

			app, err := NewApp(context.Background(), NewAppRequest{
				URL:         server.URL,
				HAAuthToken: "token",
			})
			require.NoError(t, err)
			app.RegisterEntityListeners(NewEntityListener().
				EntityIDs("light.startup").
				Call(func(*Service, State, EntityData) {}).
				RunOnStartup().
				Build())

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			runDone := make(chan error, 1)
			go func() { runDone <- app.Start(ctx) }()

			select {
			case <-requestStarted:
			case <-time.After(time.Second):
				t.Fatal("Start did not request RunOnStartup state")
			}
			require.NoError(t, stop.call(cancel, app))

			select {
			case err := <-runDone:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(time.Second):
				t.Fatal("Start did not return after cancellation")
			}
			select {
			case <-requestCanceled:
			case <-time.After(time.Second):
				t.Fatal("RunOnStartup state request was not canceled")
			}
		})
	}
}
