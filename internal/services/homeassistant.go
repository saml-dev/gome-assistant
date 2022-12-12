package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
)

type HomeAssistant struct {
	conn *websocket.Conn
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

	ws.WriteMessage(req, ha.conn, ha.ctx)
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

	ws.WriteMessage(req, ha.conn, ha.ctx)
}

func (ha *HomeAssistant) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_off"

	ws.WriteMessage(req, ha.conn, ha.ctx)
}
