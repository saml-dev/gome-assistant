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
	req := CallServiceRequest{
		Domain:  "number",
		Service: "set_value",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: map[string]any{"value": value},
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}
