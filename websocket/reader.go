package websocket

import (
	"encoding/json"
	"log/slog"
)

// Run processes incoming messages from `Conn`. It reads
// JSON-formatted messages from `conn`, partly deserializes them, and
// passes them to the subscriber that has subscribed to that message
// ID (if any). If there is an error, return the error and stop
// listening.
//
// Note that subscribers are invoked synchronously, in the same order
// as the messages arrive, and only one is run at a time. If the
// subscriber wants processing to happen in the background, it must
// spawn a goroutine itself. A subscriber is allowed to unsubscribe
// itself synchronously within the callback, in which case it is
// guaranteed not to be invoked again for subsequent messages.
func (conn *Conn) Run() error {
	for {
		bytes, err := conn.readMessage()
		if err != nil {
			return err
		}

		var msg Message
		if err := json.Unmarshal(bytes, &msg.BaseMessage); err != nil {
			slog.Warn(
				"error unmarshaling websocket message; ignoring message",
				"error", err,
				"message", string(bytes),
			)
			continue
		}
		msg.Raw = bytes

		// If a subscriber has been registered for this message ID,
		// then call it:
		if subr, ok := conn.getSubscriber(msg.ID); ok {
			subr(msg)
		}
	}
}
