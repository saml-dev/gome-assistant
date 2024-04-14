package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Switch struct {
	conn *websocket.Conn
}

func NewSwitch(conn *websocket.Conn) *Switch {
	return &Switch{
		conn: conn,
	}
}

/* Public API */

func (s Switch) TurnOn(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "switch"
	req.Service = "turn_on"
	req.Target.EntityID = entityID

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (s Switch) Toggle(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "switch"
	req.Service = "toggle"
	req.Target.EntityID = entityID

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (s Switch) TurnOff(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "switch"
	req.Service = "turn_off"
	req.Target.EntityID = entityID

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
