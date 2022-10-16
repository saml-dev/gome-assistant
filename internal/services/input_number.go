package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputNumber struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

func (ib InputNumber) Set(entityId string, value float32) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_number"
	req.Service = "set_value"
	req.ServiceData = map[string]any{"value": value}

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputNumber) Increment(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_number"
	req.Service = "increment"

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputNumber) Decrement(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_number"
	req.Service = "decrement"

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputNumber) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_number"
	req.Service = "reload"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}
