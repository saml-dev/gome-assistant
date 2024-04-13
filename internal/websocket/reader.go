package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type BaseMessage struct {
	Type    string `json:"type"`
	Id      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMsg struct {
	Type    string
	Id      int64
	Success bool
	Raw     []byte
}

// subscribe creates a new (unique) subscription number and subscribes
// `subscriber` to it.
func (conn *Conn) subscribe(subscriber Subscriber) Subscription {
	conn.subscribeMutex.Lock()
	defer conn.subscribeMutex.Unlock()

	conn.lastID++
	id := conn.lastID
	conn.subscribers[id] = subscriber
	return Subscription{
		conn: conn,
		id:   conn.lastID,
	}
}

// unsubscribe unsubscribes from `subscription`. It must be called
// exactly once for each subscription.
func (conn *Conn) unsubscribe(id int64) error {
	conn.subscribeMutex.Lock()
	defer conn.subscribeMutex.Unlock()

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

// WatchEvents subscribes to events of the given type, invoking
// `subscriber` when any such events are received. Calls to
// `subscriber` are synchronous with respect to any other received
// messages, but asynchronous with respect to writes.
func (conn *Conn) WatchEvents(eventType string, subscriber Subscriber) (Subscription, error) {
	// Make sure we're listening before events might start arriving:
	subscription := conn.subscribe(subscriber)

	e := SubEvent{
		Id:        subscription.ID(),
		Type:      "subscribe_events",
		EventType: eventType,
	}
	err := conn.WriteMessage(e)
	if err != nil {
		conn.unsubscribe(subscription.ID())
		return Subscription{}, fmt.Errorf("error writing to websocket: %w", err)
	}
	// m, _ := ReadMessage(conn, ctx)
	// log.Default().Println(string(m))
	return subscription, nil
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
			Id:      base.Id,
			Success: base.Success,
			Raw:     bytes,
		}

		if subscriber, ok := conn.getSubscriber(chanMsg.Id); ok {
			subscriber(chanMsg)
		}
	}
}
