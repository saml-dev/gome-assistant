package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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

func (ib InputNumber) Set(target ga.Target, value float32) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_number", "set_value",
		map[string]any{"value": value},
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputNumber) Increment(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_number", "increment",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputNumber) Decrement(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_number", "decrement",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputNumber) Reload() (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_number", "reload", nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
