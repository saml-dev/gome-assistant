package services

import (
	"context"
	"fmt"
	"time"

	ws "github.com/saml-dev/gome-assistant/internal/websocket"
	"nhooyr.io/websocket"
)

/* Structs */

type InputDatetime struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

func (ib InputDatetime) Set(entityId string, value time.Time) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_datetime"
	req.Service = "set_datetime"
	req.ServiceData = map[string]any{
		"timestamp": fmt.Sprint(value.Unix()),
	}

	ws.WriteMessage(req, ib.conn, ib.ctx) // TODO: this ain't working for some reason
}

func (ib InputDatetime) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_datetime"
	req.Service = "reload"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}
