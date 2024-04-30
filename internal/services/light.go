package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
func (l Light) TurnOn(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "light", "turn_on", serviceData, target,
	)
}

// Toggle a light entity.
func (l Light) Toggle(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "light", "toggle", serviceData, target,
	)
}

func (l Light) TurnOff(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "light", "turn_off", nil, target,
	)
}
