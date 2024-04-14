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

func (ib InputDatetime) Set(entityId string, value time.Time) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_datetime"
	req.Service = "set_datetime"
	req.ServiceData = map[string]any{
		"timestamp": fmt.Sprint(value.Unix()),
	}

	ib.conn.WriteMessage(req)
}

func (ib InputDatetime) Reload() {
	req := NewBaseServiceRequest(ib.conn, "")
	req.Domain = "input_datetime"
	req.Service = "reload"
	ib.conn.WriteMessage(req)
}
