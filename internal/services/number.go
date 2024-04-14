package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Number struct {
	conn *websocket.Conn
}

func NewNumber(conn *websocket.Conn) *Number {
	return &Number{
		conn: conn,
	}
}

/* Public API */

func (ib Number) SetValue(entityID string, value float32) {
	req := CallServiceRequest{}
	req.Domain = "number"
	req.Service = "set_value"
	req.Target.EntityID = entityID
	req.ServiceData = map[string]any{"value": value}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
