package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

type Event struct {
	conn *websocket.Conn
}

func NewEvent(conn *websocket.Conn) *Event {
	return &Event{
		conn: conn,
	}
}

// Fire an event
type FireEventRequest struct {
	ID        int64          `json:"id"`
	Type      string         `json:"type"` // always set to "fire_event"
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data,omitempty"`
}

/* Public API */

// Fire an event. Takes an event type and an optional map that is sent
// as `event_data`.
func (e Event) Fire(eventType string, eventData map[string]any) {
	req := FireEventRequest{
		Type: "fire_event",
	}

	req.EventType = eventType
	req.EventData = eventData

	e.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
