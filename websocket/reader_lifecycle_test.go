package websocket

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func newPipeWebsocket(t *testing.T) (*Conn, net.Conn) {
	t.Helper()

	clientPipe, serverPipe := net.Pipe()
	handshake := make(chan error, 1)
	go func() {
		request, err := httpReadRequest(serverPipe)
		if err != nil {
			handshake <- err
			return
		}

		key := request.Header.Get("Sec-WebSocket-Key")
		hash := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
		_, err = io.WriteString(serverPipe, "HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: "+base64.StdEncoding.EncodeToString(hash[:])+"\r\n"+
			"\r\n")
		handshake <- err
	}()

	websocketConn, _, err := websocket.NewClient(
		clientPipe,
		&url.URL{Scheme: "ws", Host: "example.test", Path: "/api/websocket"},
		nil,
		0,
		0,
	)
	require.NoError(t, err)
	require.NoError(t, <-handshake)

	t.Cleanup(func() {
		_ = websocketConn.Close()
		_ = serverPipe.Close()
	})
	return &Conn{conn: websocketConn}, serverPipe
}

func httpReadRequest(conn net.Conn) (*http.Request, error) {
	return http.ReadRequest(bufio.NewReader(conn))
}

func waitForDone(t *testing.T, conn *Conn) {
	t.Helper()
	select {
	case <-conn.Done():
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for websocket terminal notification")
	}
}

func TestConnRunCancellationUnblocksReadAndNotifies(t *testing.T) {
	conn, _ := newPipeWebsocket(t)
	ctx, cancel := context.WithCancel(context.Background())
	runResult := make(chan error, 1)
	go func() { runResult <- conn.Run(ctx) }()

	cancel()
	select {
	case err := <-runResult:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("cancellation did not unblock websocket read")
	}
	waitForDone(t, conn)
	require.ErrorIs(t, conn.Err(), context.Canceled)
}

func TestConnRunPreservesUnexpectedPeerClose(t *testing.T) {
	conn, peer := newPipeWebsocket(t)
	runResult := make(chan error, 1)
	go func() { runResult <- conn.Run(context.Background()) }()

	require.NoError(t, peer.Close())
	select {
	case err := <-runResult:
		require.Error(t, err)
		require.NotErrorIs(t, err, context.Canceled)
		require.ErrorIs(t, conn.Err(), err)
	case <-time.After(time.Second):
		t.Fatal("peer close did not terminate websocket read")
	}
	waitForDone(t, conn)
}

func TestConnCloseBeforeRunNotifiesWaiters(t *testing.T) {
	conn, _ := newPipeWebsocket(t)
	done := conn.Done()

	require.NoError(t, conn.Close())
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Close did not notify websocket lifecycle waiters")
	}
	require.NoError(t, conn.Err())
}

func TestConnRunLifecycleObservationsAreRaceSafe(t *testing.T) {
	conn, _ := newPipeWebsocket(t)
	ctx, cancel := context.WithCancel(context.Background())
	runResult := make(chan error, 1)
	go func() { runResult <- conn.Run(ctx) }()

	const observers = 16
	var waitGroup sync.WaitGroup
	waitGroup.Add(observers)
	for i := 0; i < observers; i++ {
		go func() {
			defer waitGroup.Done()
			for j := 0; j < 1000; j++ {
				_ = conn.Done()
				_ = conn.Err()
			}
		}()
	}
	cancel()
	err := <-runResult
	waitGroup.Wait()
	waitForDone(t, conn)
	require.ErrorIs(t, err, context.Canceled)
	require.True(t, errors.Is(conn.Err(), context.Canceled))
}
