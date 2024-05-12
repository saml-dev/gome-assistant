package websocket

import (
	"encoding/json"
	"log/slog"
)

// Start reads JSON-formatted messages from `conn`, partly
// deserializes them, and processes them. If the message ID is
// currently subscribed to, invoke the subscriber for the message. If
// there is an error reading from `conn`, log it and return.
func (conn *Conn) Start() {
	for {
		b, err := conn.readMessage()
		if err != nil {
			slog.Error("Error reading from websocket:", err)
			return
		}

		var msg Message
		if err := json.Unmarshal(b, &msg); err != nil {
			slog.Error("Error parsing JSON message from websocket:", err)
			return
		}
		// We've only deserialized part of the message, so store the
		// raw bytes as well, so that the listeners can handle them.
		msg.Raw = b

		if subscriber, ok := conn.getSubscriber(msg.ID); ok {
			subscriber(msg)
		}
	}
}
