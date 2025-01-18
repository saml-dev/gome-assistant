package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Switch struct {
	conn *ws.WebsocketWriter
}

/* Public API */

func (s Switch) TurnOn(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_on"

	return s.conn.WriteMessage(req)
}

func (s Switch) Toggle(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	return s.conn.WriteMessage(req)
}

func (s Switch) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_off"

	return s.conn.WriteMessage(req)
}
