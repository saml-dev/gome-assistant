package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputText struct {
	conn *websocket.Conn
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

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputText) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_text"
	req.Service = "reload"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}
