package websocket

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

var testWebsocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

type websocketServerResult struct {
	authMessage []byte
	serverErr   error
	closeErr    error
}

func waitForWebsocketServerResult(t *testing.T, result <-chan websocketServerResult) websocketServerResult {
	t.Helper()

	select {
	case serverResult := <-result:
		return serverResult
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for websocket server")
		return websocketServerResult{}
	}
}

func requireClientClosed(t *testing.T, serverResult websocketServerResult) {
	t.Helper()

	require.NoError(t, serverResult.serverErr)
	require.Error(t, serverResult.closeErr)
	if netError, ok := serverResult.closeErr.(net.Error); ok {
		require.False(t, netError.Timeout(), "client did not close the websocket")
	}
}

func TestNewConnRejectsUnexpectedInitialMessageAndCloses(t *testing.T) {
	serverResult := make(chan websocketServerResult, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverConn, err := testWebsocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			serverResult <- websocketServerResult{serverErr: err}
			return
		}
		defer serverConn.Close()

		if err := serverConn.WriteJSON(map[string]string{"type": "auth_ok"}); err != nil {
			serverResult <- websocketServerResult{serverErr: err}
			return
		}

		serverConn.SetReadDeadline(time.Now().Add(time.Second))
		_, message, err := serverConn.ReadMessage()
		serverResult <- websocketServerResult{authMessage: message, closeErr: err}
	}))
	defer server.Close()

	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	conn, err := NewConn(context.Background(), baseURL, "token")
	require.Error(t, err)
	require.Nil(t, conn)

	result := waitForWebsocketServerResult(t, serverResult)
	require.Empty(t, result.authMessage, "auth must not be sent after an invalid initial frame")
	requireClientClosed(t, result)
}

func TestNewConnClosesAfterAuthFailure(t *testing.T) {
	serverResult := make(chan websocketServerResult, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverConn, err := testWebsocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			serverResult <- websocketServerResult{serverErr: err}
			return
		}
		defer serverConn.Close()

		if err := serverConn.WriteJSON(map[string]string{"type": "auth_required"}); err != nil {
			serverResult <- websocketServerResult{serverErr: err}
			return
		}

		serverConn.SetReadDeadline(time.Now().Add(time.Second))
		_, authMessage, err := serverConn.ReadMessage()
		if err != nil {
			serverResult <- websocketServerResult{serverErr: err}
			return
		}

		if err := serverConn.WriteJSON(map[string]string{
			"type":    "auth_invalid",
			"message": "invalid auth token",
		}); err != nil {
			serverResult <- websocketServerResult{serverErr: err}
			return
		}

		serverConn.SetReadDeadline(time.Now().Add(time.Second))
		_, _, closeErr := serverConn.ReadMessage()
		serverResult <- websocketServerResult{authMessage: authMessage, closeErr: closeErr}
	}))
	defer server.Close()

	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	conn, err := NewConn(context.Background(), baseURL, "token")
	require.ErrorIs(t, err, ErrInvalidToken)
	require.Nil(t, conn)

	result := waitForWebsocketServerResult(t, serverResult)
	var authMessage AuthMessage
	require.NoError(t, json.Unmarshal(result.authMessage, &authMessage))
	require.Equal(t, AuthMessage{MsgType: "auth", AccessToken: "token"}, authMessage)
	requireClientClosed(t, result)
}

type countingConn struct {
	net.Conn
	closeCalls atomic.Int32
}

func (conn *countingConn) Close() error {
	conn.closeCalls.Add(1)
	return conn.Conn.Close()
}

func TestConnCloseIsIdempotentAndConcurrent(t *testing.T) {
	clientPipe, serverPipe := net.Pipe()
	defer serverPipe.Close()

	serverHandshake := make(chan error, 1)
	go func() {
		request, err := http.ReadRequest(bufio.NewReader(serverPipe))
		if err != nil {
			serverHandshake <- err
			return
		}

		key := request.Header.Get("Sec-WebSocket-Key")
		hash := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
		accept := base64.StdEncoding.EncodeToString(hash[:])
		_, err = io.WriteString(serverPipe, "HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: "+accept+"\r\n"+
			"\r\n")
		serverHandshake <- err
	}()

	trackedConn := &countingConn{Conn: clientPipe}
	websocketConn, _, err := websocket.NewClient(
		trackedConn,
		&url.URL{Scheme: "ws", Host: "example.test", Path: "/api/websocket"},
		nil,
		0,
		0,
	)
	require.NoError(t, err)
	require.NoError(t, <-serverHandshake)

	conn := &Conn{conn: websocketConn}
	const callers = 32
	errors := make(chan error, callers)
	var waitGroup sync.WaitGroup
	waitGroup.Add(callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer waitGroup.Done()
			errors <- conn.Close()
		}()
	}
	waitGroup.Wait()
	close(errors)

	for err := range errors {
		require.NoError(t, err)
	}
	require.Equal(t, int32(1), trackedConn.closeCalls.Load())
}
