package websocket

// Subscription represents a websocket-level subscription to a
// particular message ID.
type Subscription struct {
	messageID int64
}

// MessageID returns the message ID that this subscription is
// subscribed to.
func (sub Subscription) MessageID() int64 {
	return sub.messageID
}

// Subscriber is called synchronously when a message is received that
// matches its subscription's message ID.
type Subscriber func(msg ChanMsg)

// getSubscriber returns the subscriber, if any, that is subscribed to
// the specified message ID.
func (conn *Conn) getSubscriber(messageID int64) (Subscriber, bool) {
	conn.subscribersLock.RLock()
	defer conn.subscribersLock.RUnlock()

	subscriber, ok := conn.subscribers[messageID]
	return subscriber, ok
}
