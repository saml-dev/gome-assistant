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

func TestStartOnlyReconnectsWhenCallerStartsItAgain(t *testing.T) {
	var connections atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/states/zone.home" {
			_, _ = w.Write([]byte(`{"entity_id":"zone.home","state":"home","attributes":{"latitude":1,"longitude":2}}`))
			return
		}
		if r.URL.Path != "/api/websocket" {
			http.NotFound(w, r)
			return
		}

		connections.Add(1)
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
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var request struct {
			ID int64 `json:"id"`
		}
		if json.Unmarshal(message, &request) == nil {
			_ = conn.WriteJSON(map[string]any{"id": request.ID, "type": "result", "success": true, "result": nil})
		}
	}))
	t.Cleanup(server.Close)

	app := fakeApp(t, server.URL)
	err := app.Start(context.Background())
	require.Error(t, err)
	time.Sleep(20 * time.Millisecond)
	require.Equal(t, int32(1), connections.Load())

	err = app.Start(context.Background())
	require.Error(t, err)
	require.Equal(t, int32(2), connections.Load())
}
