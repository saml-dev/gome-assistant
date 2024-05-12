package services

import (
	"context"
	"fmt"
	"time"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type InputDatetime struct {
	service Service
}

func NewInputDatetime(service Service) *InputDatetime {
	return &InputDatetime{
		service: service,
	}
}

/* Public API */

func (ib InputDatetime) Set(target ga.Target, value time.Time) (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_datetime", "set_datetime",
		map[string]any{
			"timestamp": fmt.Sprint(value.Unix()),
		},
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ib InputDatetime) Reload() (any, error) {
	ctx := context.TODO()
	var result any
	err := ib.service.CallService(
		ctx, "input_datetime", "reload", nil, ga.Target{}, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
