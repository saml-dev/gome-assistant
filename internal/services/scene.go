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
	req := CallServiceRequest{
		Domain:      "scene",
		Service:     "apply",
		ServiceData: serviceData,
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Create a scene entity.
func (s Scene) Create(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "scene",
		Service: "create",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Reload the scenes.
func (s Scene) Reload() {
	req := CallServiceRequest{
		Domain:  "scene",
		Service: "reload",
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// TurnOn a scene entity.
func (s Scene) TurnOn(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "scene",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	s.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
