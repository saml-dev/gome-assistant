package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputButton struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

func (ib InputButton) Press(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_button"
	req.Service = "press"

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputButton) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_button"
	req.Service = "reload"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}
