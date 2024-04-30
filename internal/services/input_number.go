package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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
		ga.EntityTarget(entityID),
	)
}

func (ib InputNumber) Increment(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "increment",
		nil, ga.EntityTarget(entityID),
	)
}

func (ib InputNumber) Decrement(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "decrement",
		nil, ga.EntityTarget(entityID),
	)
}

func (ib InputNumber) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "reload", nil, ga.Target{},
	)
}
