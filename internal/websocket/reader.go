package websocket

import (
	"encoding/json"
	"log/slog"
)

type BaseMessage struct {
	Type    string `json:"type"`
	ID      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMsg struct {
	Type    string
	ID      int64
	Success bool
	Raw     []byte
}

// Run processes incoming messages from `Conn`. It reads
// JSON-formatted messages from `conn`, partly deserializes them, and
// passes them to the subscriber that has subscribed to that message
// ID (if any). If there is an error, return the error and stop
// listening.
//
// Note that the subscribers are invoked synchronously, in the same
// order as the messages arrived, and only one is run at a time. If
// the subscriber wants processing to happen in the background, it
// must spawn a goroutine itself.
func (conn *Conn) Run() error {
	for {
		bytes, err := conn.readMessage()
		if err != nil {
			return err
		}

		base := BaseMessage{
			// default to true for messages that don't include "success" at all
			Success: true,
		}
		_ = json.Unmarshal(bytes, &base)
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
		}

		// Result messages are sent in response to the initial subscribe request.
		// As a result, every event listener was being called on startup. This
		// check prevents that.
		if base.Type == "result" {
			return nil
		}

		chanMsg := ChanMsg{
			Type:    base.Type,
			ID:      base.ID,
			Success: base.Success,
			Raw:     bytes,
		}

		// If a subscriber has been registered for this message ID,
		// then call it, too:
		if subr, ok := conn.getSubscriber(base.ID); ok {
			subr(chanMsg)
		}
	}
}
