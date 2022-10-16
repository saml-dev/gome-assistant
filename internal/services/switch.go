package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

/* Structs */

type Switch struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

func (s Switch) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_on"

	ws.WriteMessage(req, s.conn, s.ctx)
}

func (s Switch) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	ws.WriteMessage(req, s.conn, s.ctx)
}

func (s Switch) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_off"
	ws.WriteMessage(req, s.conn, s.ctx)
}
