package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Switch struct {
	service Service
}

func NewSwitch(service Service) *Switch {
	return &Switch{
		service: service,
	}
}

/* Public API */

func (s Switch) TurnOn(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "turn_on",
		nil, EntityTarget(entityID),
	)
}

func (s Switch) Toggle(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "toggle",
		nil, EntityTarget(entityID),
	)
}

func (s Switch) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "turn_off",
		nil, EntityTarget(entityID),
	)
}
