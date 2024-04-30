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

func (s Switch) TurnOn(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "turn_on",
		nil, target,
	)
}

func (s Switch) Toggle(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "toggle",
		nil, target,
	)
}

func (s Switch) TurnOff(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "switch", "turn_off",
		nil, target,
	)
}
