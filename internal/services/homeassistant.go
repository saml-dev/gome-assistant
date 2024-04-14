package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

type HomeAssistant struct {
	conn *websocket.Conn
}

func NewHomeAssistant(conn *websocket.Conn) *HomeAssistant {
	return &HomeAssistant{
		conn: conn,
	}
}

// TurnOn a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(ha.conn, entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_on"
	req.ServiceData = serviceData

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(ha.conn, entityId)
	req.Domain = "homeassistant"
	req.Service = "toggle"
	req.ServiceData = serviceData

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ha *HomeAssistant) TurnOff(entityId string) {
	req := NewBaseServiceRequest(ha.conn, entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_off"

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
