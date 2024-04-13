package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputBoolean struct {
	conn *ws.WebsocketConn
}

/* Public API */

func (ib InputBoolean) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_on"

	ib.conn.WriteMessage(req)
}

func (ib InputBoolean) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "toggle"

	ib.conn.WriteMessage(req)
}

func (ib InputBoolean) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_off"
	ib.conn.WriteMessage(req)
}

func (ib InputBoolean) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_boolean"
	req.Service = "reload"
	ib.conn.WriteMessage(req)
}
