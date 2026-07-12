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
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var ErrInvalidToken = errors.New("invalid authentication token")

type AuthMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

type Conn struct {
	conn      *websocket.Conn
	closeOnce sync.Once
	closeErr  error

	terminalOnce sync.Once
	terminalDone chan struct{}
	terminalMu   sync.RWMutex
	terminalErr  error
	runStarted   atomic.Bool

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

	return newConn(ctx, gConn, authToken)
}

// newConn completes authentication on an already-established websocket. Keeping
// this separate from dialing also lets lifecycle tests use a port-free net.Pipe.
func newConn(ctx context.Context, gConn *websocket.Conn, authToken string) (*Conn, error) {
	conn := Conn{
		conn:         gConn,
		subscribers:  make(map[int64]Subscriber),
		terminalDone: make(chan struct{}),
	}
	closeConnOnError := true
	defer func() {
		if closeConnOnError {
			_ = conn.Close()
		}
	}()

	authComplete := atomic.Bool{}
	watcherDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			if !authComplete.Load() {
				_ = conn.Close()
			}
		case <-watcherDone:
		}
	}()
	defer close(watcherDone)

	// Read and validate the auth_required message before sending auth.
	err := conn.verifyAuthRequired(ctx)
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

	authComplete.Store(true)
	closeConnOnError = false
	return &conn, nil
}

func (conn *Conn) Close() error {
	conn.closeOnce.Do(func() {
		if conn.conn != nil {
			conn.closeErr = conn.conn.Close()
		}
		// A connection can be closed during application startup, before Run
		// begins. In that case no reader will observe the close and publish the
		// terminal state, so callers waiting on Done would otherwise block.
		if !conn.runStarted.Load() {
			conn.recordTerminal(conn.closeErr)
		}
	})
	return conn.closeErr
}

func (conn *Conn) sendAuthMessage(ctx context.Context, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := conn.conn.WriteJSON(AuthMessage{MsgType: "auth", AccessToken: token})
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return err
	}
	return nil
}

type authResponse struct {
	MsgType string `json:"type"`
	Message string `json:"message"`
}

type authRequiredMessage struct {
	MsgType string `json:"type"`
}

func (conn *Conn) verifyAuthRequired(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	msg, err := conn.readMessage()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return err
	}

	var authReq authRequiredMessage
	if err := json.Unmarshal(msg, &authReq); err != nil {
		return err
	}
	if authReq.MsgType != "auth_required" {
		return fmt.Errorf("unexpected initial websocket message type %q", authReq.MsgType)
	}

	return nil
}

func (conn *Conn) verifyAuthResponse(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	msg, err := conn.readMessage()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
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
