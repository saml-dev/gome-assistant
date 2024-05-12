package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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

func (s Switch) TurnOn(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "switch", "turn_on",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Switch) Toggle(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "switch", "toggle",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Switch) TurnOff(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "switch", "turn_off",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
