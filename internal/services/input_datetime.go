package services

import (
	"fmt"
	"time"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputDatetime struct {
	conn *ws.WebsocketConn
}

/* Public API */

func (ib InputDatetime) Set(entityId string, value time.Time) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_datetime"
	req.Service = "set_datetime"
	req.ServiceData = map[string]any{
		"timestamp": fmt.Sprint(value.Unix()),
	}

	return ib.conn.WriteMessage(req)
}

func (ib InputDatetime) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "input_datetime"
	req.Service = "reload"
	return ib.conn.WriteMessage(req)
}
