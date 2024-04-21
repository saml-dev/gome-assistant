package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputBoolean struct {
	conn *websocket.Conn
}

func NewInputBoolean(conn *websocket.Conn) *InputBoolean {
	return &InputBoolean{
		conn: conn,
	}
}

/* Public API */

func (ib InputBoolean) TurnOn(entityID string) {
	req := CallServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (ib InputBoolean) Toggle(entityID string) {
	req := CallServiceRequest{
		Domain:  "input_boolean",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (ib InputBoolean) TurnOff(entityID string) {
	req := CallServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_off",
		Target: Target{
			EntityID: entityID,
		},
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (ib InputBoolean) Reload() {
	req := CallServiceRequest{
		Domain:  "input_boolean",
		Service: "reload",
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}
