package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

type HomeAssistant struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
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

	ha.conn.WriteMessage(req, ha.ctx)
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

	ha.conn.WriteMessage(req, ha.ctx)
}

func (ha *HomeAssistant) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_off"

	ha.conn.WriteMessage(req, ha.ctx)
}
