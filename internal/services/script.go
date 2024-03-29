package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Script struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "reload"

	s.conn.WriteMessage(s.ctx, req)
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "toggle"

	s.conn.WriteMessage(s.ctx, req)
}

// Turn off a script that was created in the HA UI.
func (s Script) TurnOff() {
	req := NewBaseServiceRequest("")
	req.Domain = "script"
	req.Service = "turn_off"

	s.conn.WriteMessage(s.ctx, req)
}

// Turn on a script that was created in the HA UI.
func (s Script) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "turn_on"

	s.conn.WriteMessage(s.ctx, req)
}
