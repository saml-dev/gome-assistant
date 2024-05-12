package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Number struct {
	service Service
}

func NewNumber(service Service) *Number {
	return &Number{
		service: service,
	}
}

/* Public API */

func (ib Number) SetValue(target ga.Target, value float32) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "number", "set_value",
		map[string]any{"value": value},
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
