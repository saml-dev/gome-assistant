package websocket

type BaseMessage struct {
	Type    string `json:"type"`
	ID      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMessage struct {
	Type    string
	ID      int64
	Success bool
	Raw     []byte
}
