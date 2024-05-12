package websocket

func (conn *Conn) getSubscriber(id int64) (Subscriber, bool) {
	conn.subscribeMutex.RLock()
	defer conn.subscribeMutex.RUnlock()

	subscriber, ok := conn.subscribers[id]
	return subscriber, ok
}

// Subscriber is called synchronously when a message with the
// subscribed `id` is received.
type Subscriber func(msg Message)

// Subscription represents a websocket-level subscription to a
// particular message ID. Incoming messages with that ID will be
// forwarded to the corresponding `Subscriber`.
type Subscription struct {
	id int64
}

func (subscription Subscription) ID() int64 {
	return subscription.id
}
