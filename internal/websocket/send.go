package websocket

// Messager is called by `Send()` while holding the `writeMutex`. It
// can send one or more messages by (for each message) allocating an
// ID using `lc.NextMessageID()` then sending the message using
// `lc.SendMessage()`. The `LockedConn` should only be used while the
// callback is running.
type Messager func(lc LockedConn) error

// Send is the primary way to write a message over the websocket
// interface. Since these messages require monotonically-increasing ID
// numbers, the work from allocating a new ID number through sending
// the message has to be done under the `writeMutex`. This is done by
// passing this function a `Messager`, which is invoked while holding
// the lock.
//
// Usage:
//
//	msg := NewFooMessage{…}
//	err := conn.Send(func(lc LockedConn) error {
//		id := lc.NextMessageID()
//		// …do anything else that needs to be done with `id`…
//		msg.ID = id
//		return lc.SendMessage(msg)
//	})
func (conn *Conn) Send(msgr Messager) error {
	conn.writeLock.Lock()
	defer conn.writeLock.Unlock()

	return msgr(lockedConn{conn: conn})
}
