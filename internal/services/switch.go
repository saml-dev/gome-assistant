package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
		nil, ga.EntityTarget(entityID),
	)
}

func (s Switch) Toggle(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "toggle",
		nil, ga.EntityTarget(entityID),
	)
}

func (s Switch) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "turn_off",
		nil, ga.EntityTarget(entityID),
	)
}
