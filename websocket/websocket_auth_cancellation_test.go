package websocket

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func readWebsocketFrame(t *testing.T, conn net.Conn) (byte, error) {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(time.Second))
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return 0, err
	}
	length := int(header[1] & 0x7f)
	if length == 126 {
		extra := make([]byte, 2)
		if _, err := io.ReadFull(conn, extra); err != nil {
			return 0, err
		}
		length = int(extra[0])<<8 | int(extra[1])
	} else if length == 127 {
		extra := make([]byte, 8)
		if _, err := io.ReadFull(conn, extra); err != nil {
			return 0, err
		}
		for _, b := range extra {
			length = length<<8 | int(b)
		}
	}
	if header[1]&0x80 != 0 {
		length += 4
	}
	_, err := io.CopyN(io.Discard, conn, int64(length))
	return header[0] & 0x0f, err
}

func requirePeerClosed(t *testing.T, serverPipe net.Conn) {
	t.Helper()
	serverPipe.SetReadDeadline(time.Now().Add(time.Second))
	var b [1]byte
	_, err := serverPipe.Read(b[:])
	require.Error(t, err)
}

func writeServerTextFrame(t *testing.T, conn net.Conn, message string) {
	t.Helper()
	payload := []byte(message)
	require.Less(t, len(payload), 126)
	_, err := conn.Write(append([]byte{0x81, byte(len(payload))}, payload...))
	require.NoError(t, err)
}

func TestNewConnCancellationWhileWaitingForAuthRequired(t *testing.T) {
	pipeConn, serverPipe := newPipeWebsocket(t)

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := newConn(ctx, pipeConn.conn, "token")
		result <- err
	}()

	cancel()
	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for NewConn")
	}
	requirePeerClosed(t, serverPipe)
}

func TestNewConnCancellationWhileWaitingForAuthResponse(t *testing.T) {
	pipeConn, serverPipe := newPipeWebsocket(t)

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := newConn(ctx, pipeConn.conn, "token")
		result <- err
	}()
	writeServerTextFrame(t, serverPipe, `{"type":"auth_required"}`)

	opcode, err := readWebsocketFrame(t, serverPipe)
	require.NoError(t, err)
	require.Equal(t, byte(0x1), opcode)
	cancel()
	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for NewConn")
	}
	requirePeerClosed(t, serverPipe)
}
