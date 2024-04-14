package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputNumber struct {
	conn *websocket.Conn
}

/* Public API */

func (ib InputNumber) Set(entityId string, value float32) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_number"
	req.Service = "set_value"
	req.ServiceData = map[string]any{"value": value}

	ib.conn.WriteMessage(req)
}

func (ib InputNumber) Increment(entityId string) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_number"
	req.Service = "increment"

	ib.conn.WriteMessage(req)
}

func (ib InputNumber) Decrement(entityId string) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_number"
	req.Service = "decrement"

	ib.conn.WriteMessage(req)
}

func (ib InputNumber) Reload() {
	req := NewBaseServiceRequest(ib.conn, "")
	req.Domain = "input_number"
	req.Service = "reload"
	ib.conn.WriteMessage(req)
}
