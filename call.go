package gomeassistant

import (
	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/services"
)

func (app *App) Call(req services.BaseServiceRequest) error {
	req.RequestType = "call_service"
	req.Id = internal.GetId()
	return app.conn.WriteMessage(req)
}
