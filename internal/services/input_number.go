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
	req := CallServiceRequest{
		Domain:  "input_number",
		Service: "set_value",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: map[string]any{"value": value},
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputNumber) Increment(entityID string) {
	req := CallServiceRequest{
		Domain:  "input_number",
		Service: "increment",
		Target: Target{
			EntityID: entityID,
		},
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputNumber) Decrement(entityID string) {
	req := CallServiceRequest{
		Domain:  "input_number",
		Service: "decrement",
		Target: Target{
			EntityID: entityID,
		},
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputNumber) Reload() {
	req := CallServiceRequest{
		Domain:  "input_number",
		Service: "reload",
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
