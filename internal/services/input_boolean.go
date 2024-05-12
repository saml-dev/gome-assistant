package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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

func (ib InputBoolean) TurnOn(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_boolean", "turn_on",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputBoolean) Toggle(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_boolean", "toggle",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputBoolean) TurnOff(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_boolean", "turn_off",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputBoolean) Reload() (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_boolean", "reload", nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
