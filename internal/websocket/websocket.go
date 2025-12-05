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

func (conn *Conn) WriteMessage(msg any) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	return conn.conn.WriteJSON(msg)
}

func (conn *Conn) readMessage() ([]byte, error) {
	_, msg, err := conn.conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

func NewConn(
	ctx context.Context, baseURL *url.URL, authToken string,
) (*Conn, error) {
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
	gConn, _, err := dialer.DialContext(ctx, urlWebsockets.String(), nil)
	if err != nil {
		slog.Error("Failed to connect to websocket. Check URI\n", "url", urlWebsockets)
		return nil, err
	}

	conn := Conn{
		conn: gConn,
	}

	// Read auth_required message
	_, err = conn.readMessage()
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Send auth message
	err = conn.sendAuthMessage(ctx, authToken)
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Verify auth message was successful
	err = conn.verifyAuthResponse(ctx)
	if err != nil {
		slog.Error("Auth token is invalid. Please double check it or create a new token in your Home Assistant profile\n")
		return nil, err
	}

	return &conn, nil
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

func (conn *Conn) sendAuthMessage(ctx context.Context, token string) error {
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

func (conn *Conn) verifyAuthResponse(ctx context.Context) error {
	msg, err := conn.readMessage()
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
	// m, _ := ReadMessage(ctx, conn)
	// log.Default().Println(string(m))
}
