package websocket

import "fmt"

type MessageWriter interface {
	NextID() int64
	Subscribe(subscriber Subscriber) Subscription
	SendMessage(msg any) error
}

// Messager is called by `Send()` while holding the `writeMutex`. It
// can send a message by allocating an ID using `mw.NextID()` then
// sending it using `mw.SendMessage()`. The `MessageWriter` should
// only be used while the callback is running.
type Messager func(mw MessageWriter) error

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
//	err := conn.Send(func(mw MessageWriter) error {
//		id := mw.NextID()
//		// …do anything else that needs to be done with `id`…
//		msg.ID = id
//		return mw.SendMessage(msg)
//	})
func (conn *Conn) Send(msgr Messager) error {
	conn.writeMutex.Lock()
	defer conn.writeMutex.Unlock()

	return msgr(connMessageWriter{conn: conn})
}

// SendMessage sends the specified message over the websocket
// connection. `msg` must be JSON-serializable and have the correct
// format and a unique, monotonically-increasing ID.
func (mw connMessageWriter) SendMessage(msg any) error {
	if err := mw.conn.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("sending websocket message to server: %w", err)
	}

	return nil
}

type connMessageWriter struct {
	conn *Conn
}

func (mw connMessageWriter) NextID() int64 {
	mw.conn.lastID++
	return mw.conn.lastID
}

// Subscribe creates a new (unique) subscription ID and subscribes
// `subscriber` to it, in the sense that the subscriber will be called
// for any responses that have that ID. This doesn't actually interact
// with the server.
func (mw connMessageWriter) Subscribe(subscriber Subscriber) Subscription {
	id := mw.NextID()
	mw.conn.subscribers[id] = subscriber
	return Subscription{
		conn: mw.conn,
		id:   id,
	}
}
