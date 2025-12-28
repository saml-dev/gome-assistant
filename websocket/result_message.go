package websocket

// BaseResultMessage represents the header of a websocket message that
// holds the result of an operation.
type BaseResultMessage struct {
	BaseMessage
	Success bool `json:"success"`
}

// ResultMessage represents the full contents of a websocket message
// that holds the result of an operation.
type ResultMessage struct {
	BaseResultMessage
	Raw []byte
}
