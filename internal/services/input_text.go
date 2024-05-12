package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type InputText struct {
	service Service
}

func NewInputText(service Service) *InputText {
	return &InputText{
		service: service,
	}
}

/* Public API */

func (ib InputText) Set(target ga.Target, value string) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_text", "set_value",
		map[string]any{
			"value": value,
		},
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputText) Reload() (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_text", "reload", nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
