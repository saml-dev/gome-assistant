package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"saml.dev/gome-assistant/websocket"
)

// resultSubscriber is a helper type for handling the result message
// sent by the server in response to some kind of request. It
// subscribes itself, sends a message to the server, captures the
// first `result` message, then unsubscribes itself. The
// `ResultMessage` or error can be read using `wait()`.
type resultSubscriber struct {
	app          *App
	subscription websocket.Subscription

	once   sync.Once
	result any
	err    error
	done   chan struct{}
}

// newResultSubscriber creates a new subscriber that writes its result
// into `result`, which must be something that `json.Unmarshal()` can
// marshal into (typically a pointer).
func newResultSubscriber(app *App, result any) *resultSubscriber {
	return &resultSubscriber{
		app:    app,
		result: result,
		done:   make(chan struct{}),
	}
}

// subscribe prepares and sends `req` to `lc`, but first subscribes
// `rs.callback` to receive the result of the request.
func (rs *resultSubscriber) subscribe(
	lc websocket.LockedConn, req websocket.Request,
) error {
	rs.subscription = lc.Subscribe(rs.callback)
	req.SetID(rs.subscription.ID())
	if err := lc.SendMessage(req); err != nil {
		lc.Unsubscribe(rs.subscription)
		return fmt.Errorf("error writing to websocket: %w", err)
	}
	return nil
}

// callback receives a single "result" message, stores the result to
// `rs`, then unsubscribes. It implements `websocket.Subscriber`.
func (rs *resultSubscriber) callback(msg websocket.Message) {
	defer rs.close()
	rs.err = msg.GetResult(rs.result)
}

// wait waits for the result message to be received by `callback()`,
// then returns it to the caller.
func (rs *resultSubscriber) wait(ctx context.Context) error {
	select {
	case <-rs.done:
		return rs.err
	case <-ctx.Done():
		rs.close()
		return ctx.Err()
	}
}

func (rs *resultSubscriber) close() {
	rs.once.Do(func() {
		close(rs.done)
		err := rs.app.wsConn.Send(func(lc websocket.LockedConn) error {
			lc.Unsubscribe(rs.subscription)
			return nil
		})
		if err != nil {
			slog.Warn("Error unsubscribing", "message_id", rs.subscription.ID())
		}
	})
}
