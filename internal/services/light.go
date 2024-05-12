package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Light struct {
	service Service
}

func NewLight(service Service) *Light {
	return &Light{
		service: service,
	}
}

/* Public API */

// TurnOn a light entity.
func (l Light) TurnOn(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := l.service.CallService(
		ctx, "light", "turn_on", serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Toggle a light entity.
func (l Light) Toggle(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := l.service.CallService(
		ctx, "light", "toggle", serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (l Light) TurnOff(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := l.service.CallService(
		ctx, "light", "turn_off", nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
