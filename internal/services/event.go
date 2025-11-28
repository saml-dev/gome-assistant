package services

import (
	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/websocket"
)

type Event struct {
	conn *websocket.Conn
}

// Fire an event
type FireEventRequest struct {
	Id        int64          `json:"id"`
	Type      string         `json:"type"` // always set to "fire_event"
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data,omitempty"`
}

/* Public API */

// Fire an event. Takes an event type and an optional map that is sent
// as `event_data`.
func (e Event) Fire(eventType string, eventData ...map[string]any) error {
	req := FireEventRequest{
		Id:   internal.GetId(),
		Type: "fire_event",
	}

	req.EventType = eventType
	if len(eventData) != 0 {
		req.EventData = eventData[0]
	}

	return e.conn.WriteMessage(req)
}
