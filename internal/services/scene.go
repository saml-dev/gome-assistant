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
func (s Scene) Apply(serviceData map[string]any) {
	req := NewBaseServiceRequest(s.conn, "")
	req.Domain = "scene"
	req.Service = "apply"
	req.ServiceData = serviceData

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Create a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) Create(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "scene"
	req.Service = "create"
	req.ServiceData = serviceData

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Reload the scenes.
func (s Scene) Reload() {
	req := NewBaseServiceRequest(s.conn, "")
	req.Domain = "scene"
	req.Service = "reload"

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// TurnOn a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(s.conn, entityId)
	req.Domain = "scene"
	req.Service = "turn_on"
	req.ServiceData = serviceData

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
