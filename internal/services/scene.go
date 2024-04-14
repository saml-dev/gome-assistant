package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Scene struct {
	conn *websocket.Conn
}

func NewScene(conn *websocket.Conn) *Scene {
	return &Scene{
		conn: conn,
	}
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData ...map[string]any) {
	req := NewBaseServiceRequest(s.conn, "")
	req.Domain = "scene"
	req.Service = "apply"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	s.conn.WriteMessage(req)
}

// Create a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) Create(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "scene"
	req.Service = "create"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	s.conn.WriteMessage(req)
}

// Reload the scenes.
func (s Scene) Reload() {
	req := NewBaseServiceRequest(s.conn, "")
	req.Domain = "scene"
	req.Service = "reload"

	s.conn.WriteMessage(req)
}

// TurnOn a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "scene"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	s.conn.WriteMessage(req)
}
