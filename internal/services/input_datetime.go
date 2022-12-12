package services

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
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

	ws.WriteMessage(req, ib.conn, ib.ctx)
}

func (ib InputDatetime) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_datetime"
	req.Service = "reload"
	ws.WriteMessage(req, ib.conn, ib.ctx)
}
