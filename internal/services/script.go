package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Script struct {
	conn *websocket.Conn
}

func NewScript(conn *websocket.Conn) *Script {
	return &Script{
		conn: conn,
	}
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityID string) {
	req := CallServiceRequest{
		Domain:  "script",
		Service: "reload",
		Target: Target{
			EntityID: entityID,
		},
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityID string) {
	req := CallServiceRequest{
		Domain:  "script",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Turn off a script that was created in the HA UI.
func (s Script) TurnOff() {
	req := CallServiceRequest{
		Domain:  "script",
		Service: "turn_off",
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Turn on a script that was created in the HA UI.
func (s Script) TurnOn(entityID string) {
	req := CallServiceRequest{
		Domain:  "script",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
