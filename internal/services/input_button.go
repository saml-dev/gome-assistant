package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputButton struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

func (ib InputButton) Press(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_button"
	req.Service = "press"

	ib.conn.WriteMessage(req, ib.ctx)
}

func (ib InputButton) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_button"
	req.Service = "reload"
	ib.conn.WriteMessage(req, ib.ctx)
}
