package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

// Call invokes an RPC and processes the result as follows:
//  1. Generate a message ID.
//  2. Subscribe to that ID.
//  3. Send `req` over the websocket
//  4. Waits for a single "result" message
//  5. Unsubscribe from ID
//  6. Unmarshal the result into `result`.
//
// `msg` must be serializable to JSON. It shouldn't have its ID filled
// in yet; that will be done within this method. `result` must be
// something that `json.Unmarshal()` can deserialize into; typically,
// it is a pointer. If the result indicates a failure
// (success==false), then return that as an error.
func (app *App) Call(
	ctx context.Context, req websocket.Request, result any,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rs := newResultSubscriber(app, result)
	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		return rs.subscribe(lc, req)
	})
	if err != nil {
		return err
	}
	return rs.wait(ctx)
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
// waits for the response. The response is evaluated; if it indicates
// an error, then this method returns that error. Otherwise, the
// "result" field is stored to `result`, which must be something that
// `json.Unmarshal()` can serialize into (typically a pointer).
func (app *App) CallService(
	ctx context.Context, domain string, service string, serviceData any, target ga.Target,
	result any,
) error {
	req := CallServiceRequest{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		Domain:      domain,
		Service:     service,
		ServiceData: serviceData,
		Target:      target,
	}

	if err := app.Call(ctx, &req, result); err != nil {
		switch target {
		case ga.Target{}:
			return fmt.Errorf("calling '%s.%s': %w", domain, service, err)
		default:
			return fmt.Errorf("calling '%s.%s' for %s: %w", domain, service, target, err)
		}
	}
	return nil
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
) (websocket.ResultMessage, websocket.Subscription, error) {
	// The result of the attempt to subscribe (i.e., the first
	// message) will be sent to this channel.
	resultReceived := false
	var resultMsg websocket.ResultMessage
	var resultErr error
	done := make(chan struct{})

	var subscription websocket.Subscription

	// Receive a single "result" message, send it to `responseCh`,
	// then unsubscribe:
	dualSubscriber := func(msg websocket.Message) {
		if msg.Type == "result" {
			if resultReceived {
				slog.Warn(
					"Error: multiple responses received for one 'subscribe' request (ignored)",
				)
				return
			}
			resultReceived = true

			defer close(done)

			resultErr = json.Unmarshal(msg.Raw, &resultMsg)
			if resultErr != nil {
				return
			}
			// FIXME: turn non-success responses into errors.
			return
		}

		// Forward other responses (i.e., the events themselves) to
		// `subscriber`:
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
		return websocket.ResultMessage{}, websocket.Subscription{}, err
	}

	select {
	case <-done:
		return resultMsg, subscription, nil
	case <-ctx.Done():
		// FIXME: unsubscribe
		return websocket.ResultMessage{}, websocket.Subscription{}, ctx.Err()
	}
}
