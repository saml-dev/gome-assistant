package gomeassistant

import (
	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/websocket"
)

// CallAndForget implements [services.API.CallAndForget].
func (app *App) CallAndForget(req services.BaseServiceRequest) error {
	reqMsg := services.CallServiceMessage{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		BaseServiceRequest: req,
	}

	return app.conn.Send(
		func(lc websocket.LockedConn) error {
			reqMsg.ID = lc.NextMessageID()
			return lc.SendMessage(reqMsg)
		},
	)
}
