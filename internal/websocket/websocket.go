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

type Conn struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

func NewConnFromURI(ctx context.Context, uri string, authToken string) (*Conn, error) {
	// Init websocket connection
	dialer := websocket.DefaultDialer
	wsConn, _, err := dialer.DialContext(ctx, uri, nil)
	if err != nil {
		slog.Error("Failed to connect to websocket. Check URI\n", "uri", uri)
		return nil, err
	}

	conn := &Conn{conn: wsConn}

	// Read auth_required message
	if _, err := conn.readMessage(); err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Send auth message
	err = conn.sendAuthMessage(authToken)
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Verify auth message was successful
	err = conn.verifyAuthResponse()
	if err != nil {
		slog.Error("Auth token is invalid. Please double check it or create a new token in your Home Assistant profile\n")
		return nil, err
	}

	return conn, nil
}

func NewConn(ctx context.Context, ip, port, authToken string) (*Conn, error) {
	uri := fmt.Sprintf("ws://%s:%s/api/websocket", ip, port)
	return NewConnFromURI(ctx, uri, authToken)
}

func NewSecureConn(ctx context.Context, ip, port, authToken string) (*Conn, error) {
	uri := fmt.Sprintf("wss://%s:%s/api/websocket", ip, port)
	return NewConnFromURI(ctx, uri, authToken)
}

func (conn *Conn) WriteMessage(msg interface{}) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	err := conn.conn.WriteJSON(msg)
	if err != nil {
		return err
	}

	return nil
}

func (conn *Conn) readMessage() ([]byte, error) {
	_, msg, err := conn.conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

func (conn *Conn) sendAuthMessage(token string) error {
	err := conn.conn.WriteJSON(AuthMessage{MsgType: "auth", AccessToken: token})
	if err != nil {
		return err
	}
	return nil
}

type authResponse struct {
	MsgType string `json:"type"`
	Message string `json:"message"`
}

func (conn *Conn) verifyAuthResponse() error {
	msg, err := conn.readMessage()
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

func SubscribeToStateChangedEvents(id int64, conn *Conn) {
	SubscribeToEventType("state_changed", conn, id)
}

func SubscribeToEventType(eventType string, conn *Conn, id ...int64) {
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
