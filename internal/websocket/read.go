package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// WatchEvents subscribes to events of the given type, invoking
// `subscriber` when any such events are received. Calls to
// `subscriber` are synchronous with respect to any other received
// messages, but asynchronous with respect to writes.
func (conn *Conn) WatchEvents(eventType string, subscriber Subscriber) (Subscription, error) {
	// Make sure we're listening before events might start arriving:
	e := SubEvent{
		BaseMessage: BaseMessage{
			Type: "subscribe_events",
		},
		EventType: eventType,
	}
	var subscription Subscription
	err := conn.Send(func(lc LockedConn) error {
		subscription = lc.Subscribe(subscriber)
		e.ID = subscription.ID()
		if err := lc.SendMessage(e); err != nil {
			conn.unsubscribe(subscription.ID())
			return fmt.Errorf("error writing to websocket: %w", err)
		}
		return nil
	})
	if err != nil {
		return Subscription{}, fmt.Errorf("error writing to websocket: %w", err)
	}
	// m, _ := ReadMessage(conn, ctx)
	// log.Default().Println(string(m))
	return subscription, nil
}

// unwatchEvents unsubscribes to events with the given `subscriptionID`. This does
// not remove the subscriber.
func (conn *Conn) unwatchEvents(subscriptionID int64) error {
	e := UnsubEvent{
		BaseMessage: BaseMessage{
			Type: "unsubscribe_events",
		},
		Subscription: subscriptionID,
	}

	err := conn.Send(func(lc LockedConn) error {
		e.ID = lc.NextID()
		return lc.SendMessage(e)
	})
	if err != nil {
		return fmt.Errorf("unsubscribing from ID %d: %w", subscriptionID, err)
	}
	// m, _ := ReadMessage(conn, ctx)
	// log.Default().Println(string(m))
	return nil
}

func (conn *Conn) WatchStateChangedEvents(subscriber Subscriber) (Subscription, error) {
	return conn.WatchEvents("state_changed", subscriber)
}

// Start reads JSON-formatted messages from `conn`, partly
// deserializes them, and processes them. If the message ID is
// currently subscribed to, invoke the subscriber for the message. If
// there is an error reading from `conn`, log it and return.
func (conn *Conn) Start() {
	for {
		bytes, err := conn.readMessage()
		if err != nil {
			slog.Error("Error reading from websocket:", err)
			return
		}

		base := BaseResultMessage{
			// default to true for messages that don't include "success" at all
			Success: true,
		}
		json.Unmarshal(bytes, &base)
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
		}
		msg := Message{
			BaseMessage: base.BaseMessage,
			Raw:         bytes,
		}

		if subscriber, ok := conn.getSubscriber(msg.ID); ok {
			subscriber(msg)
		}
	}
}
