package websocket

import (
	"encoding/json"
)

// BaseMessage implements the required part of any websocket message.
// The idea is to embed this type in other message types.
type BaseMessage struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

type Request interface {
	GetID() int64
	SetID(id int64)
}

func (msg *BaseMessage) GetID() int64 {
	return msg.ID
}

func (msg *BaseMessage) SetID(id int64) {
	msg.ID = id
}

type BaseResultMessage struct {
	BaseMessage
	Success bool `json:"success"`
}

// Message holds a complete message, only partly parsed. The entire,
// original, unparsed message is available in the `Raw` field.
type Message struct {
	BaseMessage

	// Raw contains the original, full, unparsed message (including
	// fields `Type` and `ID`, which appear in `BaseMessage`).
	Raw json.RawMessage
}
