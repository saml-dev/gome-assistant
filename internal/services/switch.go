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

func (s Switch) TurnOn(entityID string) {
	req := CallServiceRequest{
		Domain:  "switch",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
	}

	s.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (s Switch) Toggle(entityID string) {
	req := CallServiceRequest{
		Domain:  "switch",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
	}

	s.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (s Switch) TurnOff(entityID string) {
	req := CallServiceRequest{
		Domain:  "switch",
		Service: "turn_off",
		Target: Target{
			EntityID: entityID,
		},
	}

	s.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}
