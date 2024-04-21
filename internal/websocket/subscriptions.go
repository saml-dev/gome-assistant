package websocket

import (
	"fmt"
)

// subscribe registers `subscriber` to be called for any messages that
// have the specified `id`. This method doesn't actually interact with
// the server; that is the caller's responsibility. This method must
// be invoked while holding the `subscribeMutex` for writing.
func (conn *Conn) subscribe(id int64, subscriber Subscriber) error {
	if _, ok := conn.subscribers[id]; ok {
		return fmt.Errorf("id %d is already subscribed to", id)
	}
	conn.subscribers[id] = subscriber
	return nil
}

func (conn *Conn) getSubscriber(id int64) (Subscriber, bool) {
	conn.subscribeMutex.RLock()
	defer conn.subscribeMutex.RUnlock()

	subscriber, ok := conn.subscribers[id]
	return subscriber, ok
}

// unsubscribe unsubscribes whatever subscriber is listening to
// `subscription`. It must be called exactly once for each
// subscription. It must be invoked while holding the `subscribeMutex`
// for writing.
func (conn *Conn) unsubscribe(id int64) error {
	if _, ok := conn.subscribers[id]; !ok {
		return fmt.Errorf("subscription ID %d wasn't active", id)
	}
	delete(conn.subscribers, id)
	return nil
}

// Subscriber is called synchronously when a message with the
// subscribed `id` is received.
type Subscriber func(msg Message)

// Subscription represents a websocket-level subscription to a
// particular message ID. Incoming messages with that ID will be
// forwarded to the corresponding `Subscriber`.
type Subscription struct {
	conn *Conn
	id   int64
}

func (subscription Subscription) ID() int64 {
	return subscription.id
}
