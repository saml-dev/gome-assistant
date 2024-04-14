package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type BaseMessage struct {
	Type    string `json:"type"`
	ID      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMsg struct {
	Type    string
	ID      int64
	Success bool
	Raw     []byte
}

// unsubscribe unsubscribes from `subscription`. It must be called
// exactly once for each subscription. It must be invoked while
// holding the `subscribeMutex` for writing.
func (conn *Conn) unsubscribe(id int64) error {
	if _, ok := conn.subscribers[id]; !ok {
		return fmt.Errorf("subscription ID %d wasn't active", id)
	}
	delete(conn.subscribers, id)
	return nil
}

func (conn *Conn) getSubscriber(id int64) (Subscriber, bool) {
	conn.subscribeMutex.RLock()
	defer conn.subscribeMutex.RUnlock()

	subscriber, ok := conn.subscribers[id]
	return subscriber, ok
}

type SubEvent struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

// WatchEvents subscribes to events of the given type, invoking
// `subscriber` when any such events are received. Calls to
// `subscriber` are synchronous with respect to any other received
// messages, but asynchronous with respect to writes.
func (conn *Conn) WatchEvents(eventType string, subscriber Subscriber) (Subscription, error) {
	// Make sure we're listening before events might start arriving:
	e := SubEvent{
		Type:      "subscribe_events",
		EventType: eventType,
	}
	var subscription Subscription
	err := conn.Send(func(mw MessageWriter) error {
		subscription = mw.Subscribe(subscriber)
		e.ID = subscription.ID()
		if err := mw.SendMessage(e); err != nil {
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

type UnsubEvent struct {
	ID           int64  `json:"id"`
	Type         string `json:"type"`
	Subscription int64  `json:"subscription"`
}

// unwatchEvents unsubscribes to events with the given `subscriptionID`. This does
// not remove the subscriber.
func (conn *Conn) unwatchEvents(subscriptionID int64) error {
	e := UnsubEvent{
		Type:         "unsubscribe_events",
		Subscription: subscriptionID,
	}

	err := conn.Send(func(mw MessageWriter) error {
		e.ID = mw.NextID()
		return mw.SendMessage(e)
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

		base := BaseMessage{
			// default to true for messages that don't include "success" at all
			Success: true,
		}
		json.Unmarshal(bytes, &base)
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
		}
		chanMsg := ChanMsg{
			Type:    base.Type,
			ID:      base.ID,
			Success: base.Success,
			Raw:     bytes,
		}

		if subscriber, ok := conn.getSubscriber(chanMsg.ID); ok {
			subscriber(chanMsg)
		}
	}
}
