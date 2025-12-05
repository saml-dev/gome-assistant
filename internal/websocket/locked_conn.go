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

// SendMessage implements [LockedConn.SendMessage].
func (lc lockedConn) SendMessage(msg any) error {
	if err := lc.conn.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("sending websocket message to server: %w", err)
	}

	return nil
}
