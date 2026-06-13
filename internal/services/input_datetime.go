package services

import (
	"context"
	"fmt"
	"time"
)

/* Structs */

type InputDatetime struct {
	api API
}

/* Public API */

func (ib InputDatetime) Set(
	ctx context.Context, target Target, value time.Time,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_datetime",
		Service: "set_datetime",
		ServiceData: map[string]any{
			"timestamp": fmt.Sprint(value.Unix()),
		},
		Target: target,
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputDatetime) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_datetime",
		Service: "reload",
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
