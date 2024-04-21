package websocket

// LockedConn represents a `Conn` object that is currently locked. It
// allows user access to operations that usually require the lock.
type LockedConn interface {
	NextID() int64
	Subscribe(subscriber Subscriber) Subscription
	Unsubscribe(subscription Subscription)
	SendMessage(msg any) error
}
