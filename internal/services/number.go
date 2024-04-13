package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Number struct {
	conn *ws.WebsocketConn
}

/* Public API */

func (ib Number) SetValue(entityId string, value float32) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "number"
	req.Service = "set_value"
	req.ServiceData = map[string]any{"value": value}

	ib.conn.WriteMessage(req)
}
