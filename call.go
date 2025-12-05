package gomeassistant

import (
	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/internal/websocket"
)

func (app *App) Call(req services.BaseServiceRequest) error {
	req.RequestType = "call_service"

	return app.conn.Send(
		func(lc websocket.LockedConn) error {
			req.Id = lc.NextMessageID()
			return lc.SendMessage(req)
		},
	)
}
