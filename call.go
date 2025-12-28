package gomeassistant

import (
	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/websocket"
)

// Call implements [services.API.Call].
func (app *App) Call(req services.BaseServiceRequest) error {
	req.Type = "call_service"

	return app.conn.Send(
		func(lc websocket.LockedConn) error {
			req.ID = lc.NextMessageID()
			return lc.SendMessage(req)
		},
	)
}
