package gomeassistant

import (
	"saml.dev/gome-assistant/internal"
)

func (app *App) FireEvent(eventType string, eventData map[string]any) error {
	req := FireEventRequest{
		Id:        internal.GetId(),
		Type:      "fire_event",
		EventType: eventType,
		EventData: eventData,
	}

	return app.conn.WriteMessage(req)
}

// Fire an event
type FireEventRequest struct {
	Id        int64          `json:"id"`
	Type      string         `json:"type"` // always set to "fire_event"
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data,omitempty"`
}
