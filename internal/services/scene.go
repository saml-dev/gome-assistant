package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Scene struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData ...map[string]any) {
	req := NewBaseServiceRequest("")
	req.Domain = "scene"
	req.Service = "apply"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, s.conn, s.ctx)
}

// Create a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) Create(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "scene"
	req.Service = "create"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, s.conn, s.ctx)
}

// Reload the scenes.
func (s Scene) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "scene"
	req.Service = "reload"

	ws.WriteMessage(req, s.conn, s.ctx)
}

// TurnOn a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "scene"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, s.conn, s.ctx)
}
