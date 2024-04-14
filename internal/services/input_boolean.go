package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputBoolean struct {
	conn *websocket.Conn
}

func NewInputBoolean(conn *websocket.Conn) *InputBoolean {
	return &InputBoolean{
		conn: conn,
	}
}

/* Public API */

func (ib InputBoolean) TurnOn(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "input_boolean"
	req.Service = "turn_on"
	req.Target.EntityID = entityID

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputBoolean) Toggle(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "input_boolean"
	req.Service = "toggle"
	req.Target.EntityID = entityID

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputBoolean) TurnOff(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "input_boolean"
	req.Service = "turn_off"
	req.Target.EntityID = entityID

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputBoolean) Reload() {
	req := CallServiceRequest{}
	req.Domain = "input_boolean"
	req.Service = "reload"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
