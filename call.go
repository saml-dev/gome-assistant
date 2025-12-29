package gomeassistant

import (
	"context"
	"sync"

	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/websocket"
)

// CallAndForget implements [services.API.CallAndForget].
func (app *App) CallAndForget(req services.BaseServiceRequest) error {
	reqMsg := services.CallServiceMessage{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		BaseServiceRequest: req,
	}

	return app.conn.Send(
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

	unsubscribe := func() {
		_ = app.conn.Send(func(lc websocket.LockedConn) error {
			lc.Unsubscribe(subscription)
			return nil
		})
	}

	handleResponse := func(msg websocket.Message) {
		once.Do(
			func() {
				responseErr = msg.GetResult(result)
				unsubscribe()
				close(done)
			},
		)
	}

	err := app.conn.Send(
		func(lc websocket.LockedConn) error {
			subscription = lc.Subscribe(handleResponse)
			reqMsg.ID = subscription.MessageID()
			return lc.SendMessage(reqMsg)
		},
	)
	if err != nil {
		return err
	}

	select {
	case <-done:
		// `handleResponse` has processed a response and set
		// `responseErr`.
	case <-ctx.Done():
		// The context has expired. Unsubscribe and return
		// `ctx.Err()`, but only if `handleResponse` hasn't just
		// racily processed a response.
		once.Do(
			func() {
				unsubscribe()
				responseErr = ctx.Err()
			},
		)
	}

	return responseErr
}
