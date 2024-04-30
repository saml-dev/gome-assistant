package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
		serviceData, ga.Target{},
	)
}

// Create a scene entity.
func (s Scene) Create(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "create",
		serviceData, target,
	)
}

// Reload the scenes.
func (s Scene) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "reload", nil, ga.Target{},
	)
}

// TurnOn a scene entity.
func (s Scene) TurnOn(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "scene", "turn_on",
		serviceData, target,
	)
}
