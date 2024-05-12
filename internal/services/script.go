package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Script struct {
	service Service
}

func NewScript(service Service) *Script {
	return &Script{
		service: service,
	}
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "script", "reload",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "script", "toggle",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Turn off a script that was created in the HA UI.
func (s Script) TurnOff() (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "script", "turn_off",
		nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Turn on a script that was created in the HA UI.
func (s Script) TurnOn(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := s.service.CallService(
		ctx, "script", "turn_on",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
