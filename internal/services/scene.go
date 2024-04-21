package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Scene struct {
	service Service
}

func NewScene(service Service) *Scene {
	return &Scene{
		service: service,
	}
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "apply",
		serviceData, Target{},
	)
}

// Create a scene entity.
func (s Scene) Create(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "create",
		serviceData, EntityTarget(entityID),
	)
}

// Reload the scenes.
func (s Scene) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "reload", nil, Target{},
	)
}

// TurnOn a scene entity.
func (s Scene) TurnOn(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "turn_on",
		serviceData, EntityTarget(entityID),
	)
}
