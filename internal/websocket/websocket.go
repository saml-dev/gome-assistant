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
)

var ErrInvalidToken = errors.New("invalid authentication token")

type Conn struct {
	writeMutex sync.Mutex
	conn       *websocket.Conn

	subscribeMutex sync.RWMutex
	subscribers    map[int64]Subscriber

	// lastID is the last message ID that has already been used. It
	// must be accessed atomically.
	lastID int64
}

// Subscriber is called synchronously when a message with the
// subscribed `id` is received.
type Subscriber func(msg Message)

type Subscription struct {
	conn *Conn
	id   int64
}

func (subscription Subscription) ID() int64 {
	return subscription.id
}

func (subscription *Subscription) Cancel() {
	if subscription.id == 0 {
		return
	}

	subscription.conn.subscribeMutex.Lock()
	defer subscription.conn.subscribeMutex.Unlock()

	subscription.conn.unsubscribe(subscription.id)

	subscription.conn.unwatchEvents(subscription.id)

	subscription.id = 0
}

func NewConnFromURI(ctx context.Context, uri string, authToken string) (*Conn, error) {
	// Init websocket connection
	dialer := websocket.DefaultDialer
	wsConn, _, err := dialer.DialContext(ctx, uri, nil)
	if err != nil {
		slog.Error("Failed to connect to websocket. Check URI\n", "uri", uri)
		return nil, err
	}

	conn := &Conn{
		conn:        wsConn,
		subscribers: make(map[int64]Subscriber),
	}

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
		slog.Error(
			"Auth token is invalid. Please double check it " +
				"or create a new token in your Home Assistant profile\n",
		)
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

type authRequest struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

func (conn *Conn) sendAuthMessage(token string) error {
	err := conn.conn.WriteJSON(authRequest{MsgType: "auth", AccessToken: token})
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
