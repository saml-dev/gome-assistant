package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Light struct {
	service Service
}

func NewLight(service Service) *Light {
	return &Light{
		service: service,
	}
}

/* Public API */

// TurnOn a light entity.
func (l Light) TurnOn(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "light", "turn_on", serviceData, EntityTarget(entityID),
	)
}

// Toggle a light entity.
func (l Light) Toggle(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "light", "toggle", serviceData, EntityTarget(entityID),
	)
}

func (l Light) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "light", "turn_off", nil, EntityTarget(entityID),
	)
}
