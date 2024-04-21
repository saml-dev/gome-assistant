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
		bytes, err := conn.readMessage()
		if err != nil {
			slog.Error("Error reading from websocket:", err)
			return
		}

		base := BaseResultMessage{
			// default to true for messages that don't include "success" at all
			Success: true,
		}
		json.Unmarshal(bytes, &base)
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
		}
		msg := Message{
			BaseMessage: base.BaseMessage,
			Raw:         bytes,
		}

		if subscriber, ok := conn.getSubscriber(msg.ID); ok {
			subscriber(msg)
		}
	}
}
