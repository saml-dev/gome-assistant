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
)

var ErrInvalidToken = errors.New("invalid authentication token")

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

type Conn struct {
	conn          *websocket.Conn
	writeLock     sync.Mutex
	lastMessageID int64

	// subscribersLock guards access to `subscribers`.
	subscribersLock sync.RWMutex

	// subscribers is a map from message ID to the subscriber that is
	// subscribed to messages with that ID.
	subscribers map[int64]Subscriber
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
		conn:        gConn,
		subscribers: make(map[int64]Subscriber),
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

func SubscribeToStateChangedEvents(conn *Conn) Subscription {
	return SubscribeToEventType("state_changed", conn)
}

func SubscribeToEventType(eventType string, conn *Conn) Subscription {
	var id int64
	err := conn.Send(
		func(lc LockedConn) error {
			id = lc.NextMessageID()
			e := SubEvent{
				Id:        id,
				Type:      "subscribe_events",
				EventType: eventType,
			}

			if err := lc.SendMessage(e); err != nil {
				return fmt.Errorf("error writing to websocket: %w", err)
			}
			// m, _ := ReadMessage(ctx, conn)
			// log.Default().Println(string(m))

			return nil
		},
	)

	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	return Subscription{id}
}
