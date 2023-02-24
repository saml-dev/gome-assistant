package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputText struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

func (ib InputText) Set(entityId string, value string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_text"
	req.Service = "set_value"
	req.ServiceData = map[string]any{
		"value": value,
	}

	ib.conn.WriteMessage(req, ib.ctx)
}

func (ib InputText) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_text"
	req.Service = "reload"
	ib.conn.WriteMessage(req, ib.ctx)
}
