package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Light struct {
	conn *websocket.Conn
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

	ws.WriteMessage(req, l.conn, l.ctx)
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

	ws.WriteMessage(req, l.conn, l.ctx)
}

func (l Light) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_off"
	ws.WriteMessage(req, l.conn, l.ctx)
}
