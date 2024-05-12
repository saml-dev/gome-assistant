package websocket

// LockedConn represents a `Conn` object that is currently locked. It
// allows user access to operations that usually require the lock.
type LockedConn interface {
	// NextID returns the next unused id to be used in a websocket
	// message. The IDs so generated must be used in order, while the
	// `LockedConn` is still active.
	NextID() int64

	// Subscribe creates a new (unique) subscription ID and subscribes
	// `subscriber` to it, in the sense that the subscriber will be
	// called for any responses that have that ID. This doesn't
	// actually interact with the server. The returned `Subscription`
	// must eventually be passed at least once to `Unsubscribe()`,
	// though `Unsubscribe()` can be called against a different
	// `LockedConn` than the one that generated it.
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
	// should be generated using `NextID()` and used in order.
	SendMessage(msg any) error
}
