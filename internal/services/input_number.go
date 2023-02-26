package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputNumber struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

func (ib InputNumber) Set(entityId string, value float32) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_number"
	req.Service = "set_value"
	req.ServiceData = map[string]any{"value": value}

	ib.conn.WriteMessage(req, ib.ctx)
}

func (ib InputNumber) Increment(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_number"
	req.Service = "increment"

	ib.conn.WriteMessage(req, ib.ctx)
}

func (ib InputNumber) Decrement(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_number"
	req.Service = "decrement"

	ib.conn.WriteMessage(req, ib.ctx)
}

func (ib InputNumber) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_number"
	req.Service = "reload"
	ib.conn.WriteMessage(req, ib.ctx)
}
