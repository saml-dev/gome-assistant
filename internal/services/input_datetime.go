package services

import (
	"fmt"
	"time"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputDatetime struct {
	conn *websocket.Conn
}

func NewInputDatetime(conn *websocket.Conn) *InputDatetime {
	return &InputDatetime{
		conn: conn,
	}
}

/* Public API */

func (ib InputDatetime) Set(entityID string, value time.Time) {
	req := CallServiceRequest{
		Domain:  "input_datetime",
		Service: "set_datetime",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: map[string]any{
			"timestamp": fmt.Sprint(value.Unix()),
		},
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputDatetime) Reload() {
	req := CallServiceRequest{
		Domain:  "input_datetime",
		Service: "reload",
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
