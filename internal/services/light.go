package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Light struct {
	conn *websocket.Conn
}

func NewLight(conn *websocket.Conn) *Light {
	return &Light{
		conn: conn,
	}
}

/* Public API */

// TurnOn a light entity.
func (l Light) TurnOn(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "light",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a light entity.
func (l Light) Toggle(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "light",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (l Light) TurnOff(entityID string) {
	req := CallServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target: Target{
			EntityID: entityID,
		},
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
