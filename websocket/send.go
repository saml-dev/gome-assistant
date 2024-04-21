package websocket

import (
	"fmt"
)

// Messager is called by `Send()` while holding the `writeMutex`. It
// can send a message by allocating an ID using `lc.NextID()` then
// sending it using `lc.SendMessage()`. The `MessageWriter` should
// only be used while the callback is running.
type Messager func(lc LockedConn) error

// Send is the primary way to write a message over the websocket
// interface. Since these messages require monotonically-increasing ID
// numbers, the work from allocating a new ID number through sending
// the message has to be done under the `writeMutex`. This is done by
// passing this function a `Messager`, which is invoked while holding
// the lock and passed the ID that it should use.
//
// Usage:
//
//	msg := NewFooMessage{…}
//	err := conn.Send(func(lc MessageWriter) error {
//		id := lc.NextID()
//		// …do anything else that needs to be done with `id`…
//		msg.ID = id
//		return lc.SendMessage(msg)
//	})
func (conn *Conn) Send(msgr Messager) error {
	conn.writeMutex.Lock()
	defer conn.writeMutex.Unlock()

	return msgr(lockedConn{conn: conn})
}

// lockedConn is a `LockedConn` view of a `Conn`, to be used
// only for a finite time when the connection is locked.
type lockedConn struct {
	conn *Conn
}

func (lc lockedConn) SendMessage(msg any) error {
	if err := lc.conn.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("sending websocket message to server: %w", err)
	}

	return nil
}

func (lc lockedConn) NextID() int64 {
	lc.conn.lastID++
	return lc.conn.lastID
}

func (lc lockedConn) Subscribe(subscriber Subscriber) Subscription {
	id := lc.NextID()
	lc.conn.subscribers[id] = subscriber
	return Subscription{
		id: id,
	}
}

func (lc lockedConn) Unsubscribe(subscription Subscription) {
	if subscription.id == 0 {
		return
	}
	delete(lc.conn.subscribers, subscription.id)
	subscription.id = 0
}
