package websocket

import (
	"encoding/json"
	"log/slog"
)

type BaseMessage struct {
	Type    string `json:"type"`
	Id      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMsg struct {
	Type    string
	Id      int64
	Success bool
	Raw     []byte
}

// ListenWebsocket reads JSON-formatted messages from `conn`, partly
// deserializes them, and passes them to `defaultSubscriber`, as well
// as the subscriber that has subscribed to that message ID (if any).
// If there is an error, return the error and stop listening.
//
// Note that the subscribers are invoked synchronously, in the same
// order as the messages arrived, and only one is run at a time. If
// the subscriber wants processing to happen in the background, it
// must spawn a goroutine itself.
func (conn *Conn) ListenWebsocket(defaultSubscriber Subscriber) error {
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
		chanMsg := ChanMsg{
			Type:    base.Type,
			Id:      base.Id,
			Success: base.Success,
			Raw:     bytes,
		}

		// Call the default subscriber in any case:
		defaultSubscriber(chanMsg)

		// If a subscriber has been registered for this message ID,
		// then call it, too:
		if subr, ok := conn.getSubscriber(base.Id); ok {
			subr(chanMsg)
		}
	}
}
