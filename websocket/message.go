package websocket

// BaseMessage implements the required part of any websocket message.
// This type can be embedded in other message types.
type BaseMessage struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

// Message holds a complete websocket message, only partly parsed. The
// entire, original, unparsed message is available in the `Raw` field.
type Message struct {
	BaseMessage

	// Raw contains the original, full, unparsed message (including
	// fields `Type` and `ID`, which also appear in `BaseMessage`).
	Raw RawMessage `json:"-"`
}
