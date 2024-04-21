package services

import (
	"context"

	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type InputBoolean struct {
	service Service
}

func NewInputBoolean(service Service) *InputBoolean {
	return &InputBoolean{
		service: service,
	}
}

/* Public API */

func (ib InputBoolean) TurnOn(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "turn_on",
		nil, EntityTarget(entityID),
	)
}

func (ib InputBoolean) Toggle(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "toggle",
		nil, EntityTarget(entityID),
	)
}

func (ib InputBoolean) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "turn_off",
		nil, EntityTarget(entityID),
	)
}

func (ib InputBoolean) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "reload", nil, Target{},
	)
}
