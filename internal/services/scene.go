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
	req := CallServiceRequest{}
	req.Domain = "scene"
	req.Service = "apply"
	req.ServiceData = serviceData

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Create a scene entity.
func (s Scene) Create(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "scene"
	req.Service = "create"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Reload the scenes.
func (s Scene) Reload() {
	req := CallServiceRequest{}
	req.Domain = "scene"
	req.Service = "reload"

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// TurnOn a scene entity.
func (s Scene) TurnOn(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "scene"
	req.Service = "turn_on"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
