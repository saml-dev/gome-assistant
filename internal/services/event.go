package services

import (
	"context"

	"saml.dev/gome-assistant/websocket"
)

type Event struct {
	service Service
}

func NewEvent(service Service) *Event {
	return &Event{
		service: service,
	}
}

// Fire an event
type FireEventRequest struct {
	websocket.BaseMessage
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data,omitempty"`
}

/* Public API */

// Fire an event. Takes an event type and an optional map that is sent
// as `event_data`.
func (e Event) Fire(eventType string, eventData map[string]any) (websocket.Message, error) {
	ctx := context.TODO()

	req := FireEventRequest{
		BaseMessage: websocket.BaseMessage{
			Type: "fire_event",
		},
	}

	req.EventType = eventType
	req.EventData = eventData

	return e.service.Call(ctx, &req)
}
