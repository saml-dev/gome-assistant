package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputBoolean struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

func (ib InputBoolean) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_on"

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputBoolean) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "toggle"

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputBoolean) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_off"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputBoolean) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_boolean"
	req.Service = "reload"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}
