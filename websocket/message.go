package websocket

import (
	"encoding/json"
)

// RawMessage is like `json.RawMessage`, except that its `String()`
// method converts it directly to a string.
type RawMessage json.RawMessage

// UnmarshalJSON delegates to `json.RawMessage`. (The method has a
// pointer receiver, so we have to implement it explicitly.)
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(m).UnmarshalJSON(data)
}

func (rm RawMessage) String() string {
	return string(rm)
}

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
	// fields `Type` and `ID`, which also appear in `BaseMessage`).
	Raw RawMessage `json:"-"`
}
