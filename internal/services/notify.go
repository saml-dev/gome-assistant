package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
	"saml.dev/gome-assistant/types"
)

type Notify struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

// Send a notification. Takes a types.NotifyRequest.
func (ha *Notify) Notify(reqData types.NotifyRequest) {
	req := NewBaseServiceRequest("")
	req.Domain = "notify"
	req.Service = reqData.ServiceName

	serviceData := map[string]any{}
	serviceData["message"] = reqData.Message
	serviceData["title"] = reqData.Title
	if reqData.Data != nil {
		serviceData["data"] = reqData.Data
	}

	req.ServiceData = serviceData
	ha.conn.WriteMessage(req, ha.ctx)
}
