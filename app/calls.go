package app

import (
	"context"
	"fmt"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

// Call invokes an RPC service corresponding to `req` via websockets
// and waits for and returns a single `result`. `msg` must be
// serializable to JSON. It shouldn't have its ID filled in yet; that
// will be done within this method. The response is not analyzed at
// all, even to check for errors.
func (app *App) Call(
	ctx context.Context, req websocket.Request,
) (websocket.Message, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseCh := make(chan websocket.Message, 1)

	var subscription websocket.Subscription

	// Receive a single message, sent it to `responseCh`, then
	// unsubscribe:
	subscriber := func(msg websocket.Message) {
		defer close(responseCh)
		responseCh <- msg
		_ = app.wsConn.Send(func(lc websocket.LockedConn) error {
			lc.Unsubscribe(subscription)
			return nil
		})
	}

	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		subscription = lc.Subscribe(subscriber)
		req.SetID(subscription.ID())
		if err := lc.SendMessage(req); err != nil {
			lc.Unsubscribe(subscription)
			return fmt.Errorf("error writing to websocket: %w", err)
		}
		return nil
	})

	if err != nil {
		return websocket.Message{}, err
	}

	select {
	case response := <-responseCh:
		return response, nil
	case <-ctx.Done():
		return websocket.Message{}, ctx.Err()
	}
}

type CallServiceRequest struct {
	websocket.BaseMessage
	Domain  string `json:"domain"`
	Service string `json:"service"`

	// ServiceData must be serializable to a JSON object.
	ServiceData any `json:"service_data,omitempty"`

	Target ga.Target `json:"target,omitempty"`
}

// CallService invokes a service using a `call_service` message, then
// waits for and returns the response.
//
// FIXME: can the response be parsed into a result-style message?
func (app *App) CallService(
	ctx context.Context, domain string, service string, serviceData any, target ga.Target,
) (websocket.Message, error) {
	req := CallServiceRequest{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		Domain:      domain,
		Service:     service,
		ServiceData: serviceData,
		Target:      target,
	}

	return app.Call(ctx, &req)
}

// Subscribe subscribes to some events via `req`, waits for a single
// response, and then leaves `subscriber` subscribed to the events. If
// this method returns without an error, `subscriber` must eventually
// be unsubscribed. `ctx` covers the subscription and the wait for the
// first answer, but not the forwarding of subsequent events or
// unsubscribing.
//
// FIXME: should this subscriber and subscription be specialized to
// event messages?
//
// FIXME: should the result be examined? If the subscription request
// failed, then we could fail more generally instead of leaving the
// cleanup to the caller.
func (app *App) Subscribe(
	ctx context.Context, req websocket.Request, subscriber websocket.Subscriber,
) (websocket.Message, websocket.Subscription, error) {
	// The result of the attempt to subscribe (i.e., the first
	// message) will be sent to this channel.
	resultReceived := false
	resultCh := make(chan websocket.Message, 1)

	var subscription websocket.Subscription

	// Receive a single message, sent it to `responseCh`, then
	// unsubscribe:
	dualSubscriber := func(msg websocket.Message) {
		if !resultReceived {
			// This is the first message. We send it to the channel so
			// that it can be returned from the outer function.
			defer close(resultCh)
			resultCh <- msg
			resultReceived = true
			return
		}

		// The result has already been processed. Subsequent events
		// get forwarded to `subscriber`:
		subscriber(msg)
	}

	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		subscription = lc.Subscribe(dualSubscriber)
		req.SetID(subscription.ID())
		if err := lc.SendMessage(req); err != nil {
			lc.Unsubscribe(subscription)
			return fmt.Errorf("error writing to websocket: %w", err)
		}
		return nil
	})

	if err != nil {
		return websocket.Message{}, websocket.Subscription{}, err
	}

	select {
	case response := <-resultCh:
		return response, subscription, nil
	case <-ctx.Done():
		return websocket.Message{}, websocket.Subscription{}, ctx.Err()
	}
}
