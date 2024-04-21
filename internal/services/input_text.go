package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputText struct {
	conn *websocket.Conn
}

func NewInputText(conn *websocket.Conn) *InputText {
	return &InputText{
		conn: conn,
	}
}

/* Public API */

func (ib InputText) Set(entityID string, value string) {
	req := CallServiceRequest{
		Domain:  "input_text",
		Service: "set_value",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: map[string]any{
			"value": value,
		},
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (ib InputText) Reload() {
	req := CallServiceRequest{
		Domain:  "input_text",
		Service: "reload",
	}

	ib.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}
