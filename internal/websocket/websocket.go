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
	"sync"

	"github.com/gorilla/websocket"

	"saml.dev/gome-assistant/internal"
)

var ErrInvalidToken = errors.New("invalid authentication token")

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

type WebsocketConn struct {
	Conn  *websocket.Conn
	mutex sync.Mutex
}

func (w *WebsocketConn) WriteMessage(msg interface{}) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	err := w.Conn.WriteJSON(msg)
	if err != nil {
		return err
	}

	return nil
}

func ReadMessage(conn *websocket.Conn) ([]byte, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

func SetupConnection(ctx context.Context, ip, port, authToken string) (*websocket.Conn, error) {
	uri := fmt.Sprintf("ws://%s:%s/api/websocket", ip, port)
	return ConnectionFromUri(ctx, uri, authToken)
}

func SetupSecureConnection(ctx context.Context, ip, port, authToken string) (*websocket.Conn, error) {
	uri := fmt.Sprintf("wss://%s:%s/api/websocket", ip, port)
	return ConnectionFromUri(ctx, uri, authToken)
}

func ConnectionFromUri(ctx context.Context, uri, authToken string) (*websocket.Conn, error) {
	// Init websocket connection
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, uri, nil)
	if err != nil {
		slog.Error("Failed to connect to websocket. Check URI\n", "uri", uri)
		return nil, err
	}

	// Read auth_required message
	_, err = ReadMessage(conn)
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Send auth message
	err = SendAuthMessage(conn, authToken)
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Verify auth message was successful
	err = VerifyAuthResponse(conn)
	if err != nil {
		slog.Error("Auth token is invalid. Please double check it or create a new token in your Home Assistant profile\n")
		return nil, err
	}

	return conn, nil
}

func SendAuthMessage(conn *websocket.Conn, token string) error {
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

func VerifyAuthResponse(conn *websocket.Conn) error {
	msg, err := ReadMessage(conn)
	if err != nil {
		return err
	}

	var authResp authResponse
	json.Unmarshal(msg, &authResp)
	// log.Println(authResp.MsgType)
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

func SubscribeToStateChangedEvents(id int64, conn *WebsocketConn) {
	SubscribeToEventType("state_changed", conn, id)
}

func SubscribeToEventType(eventType string, conn *WebsocketConn, id ...int64) {
	var finalId int64
	if len(id) == 0 {
		finalId = internal.GetId()
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
	// m, _ := ReadMessage(conn, ctx)
	// log.Default().Println(string(m))
}
