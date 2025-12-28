package gomeassistant

import "saml.dev/gome-assistant/websocket"

// FireEvent implements [services.API.FireEvent].
func (app *App) FireEvent(eventType string, eventData map[string]any) error {
	return app.conn.Send(
		func(lc websocket.LockedConn) error {
			req := FireEventRequest{
				ID:        lc.NextMessageID(),
				Type:      "fire_event",
				EventType: eventType,
				EventData: eventData,
			}

			return lc.SendMessage(req)
		},
	)
}

// Fire an event
type FireEventRequest struct {
	ID        int64          `json:"id"`
	Type      string         `json:"type"` // always set to "fire_event"
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data,omitempty"`
}
