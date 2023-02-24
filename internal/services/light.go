package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Light struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// TurnOn a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) TurnOn(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	l.conn.WriteMessage(req, l.ctx)
}

// Toggle a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) Toggle(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "toggle"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	l.conn.WriteMessage(req, l.ctx)
}

func (l Light) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_off"
	l.conn.WriteMessage(req, l.ctx)
}
