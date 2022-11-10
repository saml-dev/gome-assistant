package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

/* Structs */

type Number struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

func (ib Number) SetValue(entityId string, value int) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "number"
	req.Service = "set_value"
	req.ServiceData = map[string]any{"value": value}

	ws.WriteMessage(req, ib.conn, ib.ctx)
}
