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

// TurnOn a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "homeassistant"
	req.Service = "turn_on"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "homeassistant"
	req.Service = "toggle"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ha *HomeAssistant) TurnOff(entityID string) {
	req := CallServiceRequest{}
	req.Domain = "homeassistant"
	req.Service = "turn_off"
	req.Target.EntityID = entityID

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
