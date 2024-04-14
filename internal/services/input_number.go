package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputNumber struct {
	conn *websocket.Conn
}

func NewInputNumber(conn *websocket.Conn) *InputNumber {
	return &InputNumber{
		conn: conn,
	}
}

/* Public API */

func (ib InputNumber) Set(entityID string, value float32) {
	req := CallServiceRequest{}
	req.Domain = "input_number"
	req.Service = "set_value"
	req.Target.EntityID = entityID
	req.ServiceData = map[string]any{"value": value}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputNumber) Increment(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "input_number"
	req.Service = "increment"
	req.Target.EntityID = entityID

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputNumber) Decrement(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "input_number"
	req.Service = "decrement"
	req.Target.EntityID = entityID

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputNumber) Reload() {
	req := CallServiceRequest{}
	req.Domain = "input_number"
	req.Service = "reload"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
