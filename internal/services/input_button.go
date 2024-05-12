package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type InputButton struct {
	service Service
}

func NewInputButton(service Service) *InputButton {
	return &InputButton{
		service: service,
	}
}

/* Public API */

func (ib InputButton) Press(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_button", "press",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputButton) Reload() (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_button", "reload", nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
