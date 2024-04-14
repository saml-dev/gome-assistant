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

	s.conn.WriteMessage(req)
}

func (s Switch) Toggle(entityId string) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	s.conn.WriteMessage(req)
}

func (s Switch) TurnOff(entityId string) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "switch"
	req.Service = "turn_off"
	s.conn.WriteMessage(req)
}
