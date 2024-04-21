package services

import (
	"context"

	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type InputNumber struct {
	service Service
}

func NewInputNumber(service Service) *InputNumber {
	return &InputNumber{
		service: service,
	}
}

/* Public API */

func (ib InputNumber) Set(entityID string, value float32) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "set_value",
		map[string]any{"value": value},
		EntityTarget(entityID),
	)
}

func (ib InputNumber) Increment(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "increment",
		nil, EntityTarget(entityID),
	)
}

func (ib InputNumber) Decrement(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "decrement",
		nil, EntityTarget(entityID),
	)
}

func (ib InputNumber) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "reload", nil, Target{},
	)
}
