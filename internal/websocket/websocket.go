// Package websocket is used to interact with the Home Assistant
// websocket API. All HA interaction is done via websocket
// except for cases explicitly called out in http package
// documentation.
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	i "saml.dev/gome-assistant/internal"
)

var ErrInvalidToken = errors.New("invalid authentication token")

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

type WebsocketWriter struct {
	Conn  *websocket.Conn
	mutex sync.Mutex
}

func (w *WebsocketWriter) WriteMessage(msg any) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.Conn.WriteJSON(msg)
}

func ReadMessage(conn *websocket.Conn) ([]byte, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

func ConnectionFromUri(baseURL *url.URL, authToken string) (*websocket.Conn, context.Context, context.CancelFunc, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*3)

	// Shallow copy the URL to avoid modifying the original
	urlWebsockets := *baseURL
	urlWebsockets.Path = "/api/websocket"
	if baseURL.Scheme == "http" {
		urlWebsockets.Scheme = "ws"
	}
	if baseURL.Scheme == "https" {
		urlWebsockets.Scheme = "wss"
	}

	// Init websocket connection
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, urlWebsockets.String(), nil)
	if err != nil {
		ctxCancel()
		slog.Error("Failed to connect to websocket. Check URI\n", "url", urlWebsockets)
		return nil, nil, nil, err
	}

	// Read auth_required message
	_, err = ReadMessage(conn)
	if err != nil {
		ctxCancel()
		slog.Error("Unknown error creating websocket client\n")
		return nil, nil, nil, err
	}

	// Send auth message
	err = SendAuthMessage(ctx, conn, authToken)
	if err != nil {
		ctxCancel()
		slog.Error("Unknown error creating websocket client\n")
		return nil, nil, nil, err
	}

	// Verify auth message was successful
	err = VerifyAuthResponse(ctx, conn)
	if err != nil {
		ctxCancel()
		slog.Error("Auth token is invalid. Please double check it or create a new token in your Home Assistant profile\n")
		return nil, nil, nil, err
	}

	return conn, ctx, ctxCancel, nil
}

func SendAuthMessage(ctx context.Context, conn *websocket.Conn, token string) error {
	err := conn.WriteJSON(AuthMessage{MsgType: "auth", AccessToken: token})
	if err != nil {
		return err
	}
	return nil
}

type authResponse struct {
	MsgType string `json:"type"`
	Message string `json:"message"`
}

func VerifyAuthResponse(ctx context.Context, conn *websocket.Conn) error {
	msg, err := ReadMessage(conn)
	if err != nil {
		return err
	}

	var authResp authResponse
	err = json.Unmarshal(msg, &authResp)
	if err != nil {
		return err
	}
	if authResp.MsgType != "auth_ok" {
		return ErrInvalidToken
	}

	return nil
}

type SubEvent struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

func SubscribeToStateChangedEvents(ctx context.Context, id int64, conn *WebsocketWriter) {
	SubscribeToEventType(ctx, "state_changed", conn, id)
}

func SubscribeToEventType(ctx context.Context, eventType string, conn *WebsocketWriter, id ...int64) {
	var finalId int64
	if len(id) == 0 {
		finalId = i.GetId()
	} else {
		finalId = id[0]
	}
	e := SubEvent{
		Id:        finalId,
		Type:      "subscribe_events",
		EventType: eventType,
	}
	err := conn.WriteMessage(e)
	if err != nil {
		wrappedErr := fmt.Errorf("error writing to websocket: %w", err)
		slog.Error(wrappedErr.Error())
		panic(wrappedErr)
	}
	// m, _ := ReadMessage(ctx, conn)
	// log.Default().Println(string(m))
}
