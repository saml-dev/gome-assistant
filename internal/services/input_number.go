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

func (ib InputNumber) Set(target ga.Target, value float32) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "set_value",
		map[string]any{"value": value},
		target,
	)
}

func (ib InputNumber) Increment(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "increment",
		nil, target,
	)
}

func (ib InputNumber) Decrement(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "decrement",
		nil, target,
	)
}

func (ib InputNumber) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_number", "reload", nil, ga.Target{},
	)
}
