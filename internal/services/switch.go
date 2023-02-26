package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Switch struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

func (s Switch) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_on"

	s.conn.WriteMessage(req, s.ctx)
}

func (s Switch) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	s.conn.WriteMessage(req, s.ctx)
}

func (s Switch) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_off"
	s.conn.WriteMessage(req, s.ctx)
}
