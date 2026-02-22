package message

type SubscribeEventRequest struct {
	BaseMessage
	EventType string `json:"event_type"`
}
