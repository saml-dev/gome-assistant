package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
	"github.com/saml-dev/gome-assistant/types"
)

type Notify struct {
	conn *websocket.Conn
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
	ws.WriteMessage(req, ha.conn, ha.ctx)
}
