package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Script struct {
	conn *websocket.Conn
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "reload"

	s.conn.WriteMessage(req)
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "toggle"

	s.conn.WriteMessage(req)
}

// Turn off a script that was created in the HA UI.
func (s Script) TurnOff() {
	req := NewBaseServiceRequest("")
	req.Domain = "script"
	req.Service = "turn_off"

	s.conn.WriteMessage(req)
}

// Turn on a script that was created in the HA UI.
func (s Script) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "turn_on"

	s.conn.WriteMessage(req)
}
