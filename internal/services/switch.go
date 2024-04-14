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

func (s Switch) TurnOn(entityId string) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "switch"
	req.Service = "turn_on"

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (s Switch) Toggle(entityId string) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (s Switch) TurnOff(entityId string) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "switch"
	req.Service = "turn_off"

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
