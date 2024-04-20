package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputButton struct {
	conn *websocket.Conn
}

func NewInputButton(conn *websocket.Conn) *InputButton {
	return &InputButton{
		conn: conn,
	}
}

/* Public API */

func (ib InputButton) Press(entityID string) {
	req := CallServiceRequest{
		Domain:  "input_button",
		Service: "press",
		Target: Target{
			EntityID: entityID,
		},
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputButton) Reload() {
	req := CallServiceRequest{
		Domain:  "input_button",
		Service: "reload",
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
