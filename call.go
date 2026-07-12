package gomeassistant

import (
	"context"
	"sync"

	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/websocket"
)

// CallAndForget implements [services.API.CallAndForget].
func (app *App) CallAndForget(req services.BaseServiceRequest) error {
	conn, err := app.activeConn()
	if err != nil {
		return err
	}

	reqMsg := services.CallServiceMessage{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		BaseServiceRequest: req,
	}

	return conn.Send(
		func(lc websocket.LockedConn) error {
			reqMsg.ID = lc.NextMessageID()
			return lc.SendMessage(reqMsg)
		},
	)
}

// Call implements [services.API.Call].
func (app *App) Call(
	ctx context.Context, req services.BaseServiceRequest, result any,
) error {
	if ctx == nil {
		return ErrInvalidArgs
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	conn, err := app.activeConn()
	if err != nil {
		return err
	}
	select {
	case <-conn.Done():
		if err := conn.Err(); err != nil {
			return err
		}
		return ErrConnectionClosed
	default:
	}

	// Call works as follows:
	//  1. Generate a message ID.
	//  2. Subscribe to that ID.
	//  3. Send a `CallServiceMessage` containing `req` over the websocket.
	//  4. Wait for a single "result" message.
	//  5. Unsubscribe from ID.
	//  6. Unmarshal the "result" part of the response into `result`.

	reqMsg := services.CallServiceMessage{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		BaseServiceRequest: req,
	}

	// once ensures that exactly one of the following occurs:
	//	* a single response is handled and then the handler
	//	 unsubscribes itself; or
	//	* (if `ctx` expires) the handler is unsubscribed if and only
	//	 if no response has been handled.
	var once sync.Once

	// responseErr is set either to the error in the response message,
	// or to `ctx.Err()`.
	var responseErr error

	// done is closed once a response has been processed.
	done := make(chan struct{})

	var subscription websocket.Subscription

	handleResponse := func(msg websocket.Message) {
		once.Do(func() {
			responseErr = msg.GetResult(result)
			_ = conn.Send(func(lc websocket.LockedConn) error {
				lc.Unsubscribe(subscription)
				return nil
			})
			close(done)
		})
	}

	err = conn.Send(
		func(lc websocket.LockedConn) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			select {
			case <-conn.Done():
				if err := conn.Err(); err != nil {
					return err
				}
				return ErrConnectionClosed
			default:
			}
			subscription = lc.Subscribe(handleResponse)
			reqMsg.ID = subscription.MessageID()
			if err := ctx.Err(); err != nil {
				lc.Unsubscribe(subscription)
				return err
			}
			select {
			case <-conn.Done():
				lc.Unsubscribe(subscription)
				if err := conn.Err(); err != nil {
					return err
				}
				return ErrConnectionClosed
			default:
			}
			if err := lc.SendMessage(reqMsg); err != nil {
				lc.Unsubscribe(subscription)
				return err
			}
			return nil
		},
	)
	if err != nil {
		// The send callback may have subscribed before failing. This is safe
		// even when it already unsubscribed, and covers all send failures.
		_ = conn.Send(func(lc websocket.LockedConn) error {
			lc.Unsubscribe(subscription)
			return nil
		})
		return err
	}

	select {
	case <-done:
		// `handleResponse` has processed a response and set
		// `responseErr`.
	case <-ctx.Done():
		once.Do(func() {
			responseErr = ctx.Err()
			_ = conn.Send(func(lc websocket.LockedConn) error {
				lc.Unsubscribe(subscription)
				return nil
			})
			close(done)
		})
	case <-conn.Done():
		once.Do(func() {
			responseErr = conn.Err()
			if responseErr == nil {
				responseErr = ErrConnectionClosed
			}
			_ = conn.Send(func(lc websocket.LockedConn) error {
				lc.Unsubscribe(subscription)
				return nil
			})
			close(done)
		})
	}

	return responseErr
}
