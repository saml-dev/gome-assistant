package websocket

import "fmt"

// LockedConn represents a `Conn` that is currently locked for
// writing. It is created within [Conn.Send] for the use of that
// method's callback function.
type LockedConn interface {
	// NextMessageID returns the next unused id to be used in a websocket
	// message. The IDs so generated must be used in order, while the
	// `LockedConn` is still active.
	NextMessageID() int64

	// Subscribe allocates a new message ID and subscribes
	// `subscriber` to it, in the sense that the subscriber will be
	// called for any incoming messages that have that ID. This
	// doesn't actually interact with the server. Typically the next
	// step would be to send a message with its message ID set to
	// `Subscription.ID()`.
	//
	// The returned `Subscription` must eventually be passed at least
	// once to `Unsubscribe()`, though `Unsubscribe()` can be called
	// against a different `LockedConn` than the one that generated
	// it.
	Subscribe(subscriber Subscriber) Subscription

	// Unsubscribe terminates `subscription` at the websocket level;
	// i.e., no more incoming messages will be forwarded to the
	// corresponding `Subscriber`. Note that this does not interact
	// with the server; it is the caller's responsibility to send it
	// an "unsubscribe" command if necessary.
	Unsubscribe(subscription Subscription)

	// SendMessage sends the specified message over the websocket
	// connection. `msg` must be JSON-serializable and have the
	// correct format and a unique, monotonically-increasing ID, which
	// should be generated using `NextMessageID()` and used in order.
	SendMessage(msg any) error
}

// lockedConn is a `LockedConn` view of a `Conn`, to be used
// only for a finite time when the connection is locked.
type lockedConn struct {
	conn *Conn
}

// NextMessageID implements [LockedConn.NextMessageID].
func (lc lockedConn) NextMessageID() int64 {
	lc.conn.lastMessageID++
	return lc.conn.lastMessageID
}

// Subscribe implements [LockedConn.Subscribe].
func (lc lockedConn) Subscribe(subscriber Subscriber) Subscription {
	lc.conn.subscribersLock.Lock()
	defer lc.conn.subscribersLock.Unlock()

	id := lc.NextMessageID()
	lc.conn.subscribers[id] = subscriber
	return Subscription{
		messageID: id,
	}
}

// Unsubscribe implements [LockedConn.Unsubscribe].
func (lc lockedConn) Unsubscribe(subscription Subscription) {
	if subscription.messageID == 0 {
		return
	}

	lc.conn.subscribersLock.Lock()
	defer lc.conn.subscribersLock.Unlock()

	delete(lc.conn.subscribers, subscription.messageID)
	subscription.messageID = 0
}

// SendMessage implements [LockedConn.SendMessage].
func (lc lockedConn) SendMessage(msg any) error {
	if err := lc.conn.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("sending websocket message to server: %w", err)
	}

	return nil
}
