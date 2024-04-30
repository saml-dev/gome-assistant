package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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

func (ib InputBoolean) TurnOn(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "turn_on",
		nil, target,
	)
}

func (ib InputBoolean) Toggle(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "toggle",
		nil, target,
	)
}

func (ib InputBoolean) TurnOff(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "turn_off",
		nil, target,
	)
}

func (ib InputBoolean) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_boolean", "reload", nil, ga.Target{},
	)
}
