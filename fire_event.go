package gomeassistant

import (
	"saml.dev/gome-assistant/message"
	"saml.dev/gome-assistant/websocket"
)

// FireEvent implements [services.API.FireEvent].
func (app *App) FireEvent(eventType string, eventData any) error {
	return app.conn.Send(
		func(lc websocket.LockedConn) error {
			req := message.FireEventRequest{
				BaseMessage: message.BaseMessage{
					ID:   lc.NextMessageID(),
					Type: "fire_event",
				},
				EventType: eventType,
				EventData: eventData,
			}

			// FIXME: wait for result to make sure that the event was
			// fired successfully.
			return lc.SendMessage(req)
		},
	)
}
