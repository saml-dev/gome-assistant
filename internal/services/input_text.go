package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputText struct {
	conn *websocket.Conn
}

/* Public API */

func (ib InputText) Set(entityId string, value string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_text"
	req.Service = "set_value"
	req.ServiceData = map[string]any{
		"value": value,
	}

	return ib.conn.WriteMessage(req)
}

func (ib InputText) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "input_text"
	req.Service = "reload"
	return ib.conn.WriteMessage(req)
}
