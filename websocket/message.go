package websocket

// BaseMessage implements the required part of any websocket message.
// This type can be embedded in other message types.
type BaseMessage struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

type ChanMessage struct {
	Type    string
	ID      int64
	Success bool
	Raw     []byte
}
