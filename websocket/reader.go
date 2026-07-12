package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
)

// Run processes incoming messages from `Conn`. It reads
// JSON-formatted messages from `conn`, partly deserializes them, and
// passes them to the subscriber that has subscribed to that message
// ID (if any). If there is an error, return the error and stop
// listening.
//
// Note that subscribers are invoked synchronously, in the same order
// as the messages arrive, and only one is run at a time. If the
// subscriber wants processing to happen in the background, it must
// spawn a goroutine itself. A subscriber is allowed to unsubscribe
// itself synchronously within the callback, in which case it is
// guaranteed not to be invoked again for subsequent messages.
func (conn *Conn) Run(ctx context.Context) error {
	conn.runStarted.Store(true)

	runFinished := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			// Closing the underlying connection unblocks ReadMessage.
			_ = conn.Close()
		case <-runFinished:
		}
	}()
	defer close(runFinished)

	for {
		bytes, err := conn.readMessage()
		if err != nil {
			if ctx.Err() != nil {
				err = ctx.Err()
			}
			conn.recordTerminal(err)
			return err
		}

		var msg Message
		if err := json.Unmarshal(bytes, &msg.BaseMessage); err != nil {
			slog.Warn(
				"error unmarshaling websocket message; ignoring message",
				"error", err,
				"message", string(bytes),
			)
			continue
		}
		msg.Raw = bytes

		// If a subscriber has been registered for this message ID,
		// then call it:
		if subr, ok := conn.getSubscriber(msg.ID); ok {
			subr(msg)
		}
	}
}

func (conn *Conn) recordTerminal(err error) {
	conn.terminalOnce.Do(func() {
		conn.terminalMu.Lock()
		if conn.terminalDone == nil {
			conn.terminalDone = make(chan struct{})
		}
		conn.terminalErr = err
		close(conn.terminalDone)
		conn.terminalMu.Unlock()
	})
}

// Done returns a channel closed when Run reaches its terminal state.
func (conn *Conn) Done() <-chan struct{} {
	conn.terminalMu.Lock()
	defer conn.terminalMu.Unlock()
	if conn.terminalDone == nil {
		conn.terminalDone = make(chan struct{})
	}
	return conn.terminalDone
}

// Err returns the error recorded when Run reached its terminal state.
func (conn *Conn) Err() error {
	conn.terminalMu.RLock()
	defer conn.terminalMu.RUnlock()
	return conn.terminalErr
}
