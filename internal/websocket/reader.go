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
// deserializes them, and passes them to `subscriber`. If there is an
// error, return the error.
func (conn *Conn) ListenWebsocket(subscriber Subscriber) error {
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

		subscriber(chanMsg)
	}
}
