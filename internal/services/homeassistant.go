package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

type HomeAssistant struct {
	conn *ws.WebsocketWriter
}

// TurnOn a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ha.conn.WriteMessage(req)
}

// Toggle a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "toggle"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ha.conn.WriteMessage(req)
}

func (ha *HomeAssistant) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_off"

	ha.conn.WriteMessage(req)
}
