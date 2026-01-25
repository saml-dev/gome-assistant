package message

// FireEventRequest requests that Home Assistant fire an event.
type FireEventRequest struct {
	BaseMessage
	EventType string `json:"event_type"`
	EventData any    `json:"event_data,omitempty"`
}
