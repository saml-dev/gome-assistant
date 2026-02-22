package services

import (
	"context"
	"fmt"
	"time"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type InputDatetime struct {
	api API
}

/* Public API */

func (ib InputDatetime) Set(
	ctx context.Context, entityID string, value time.Time,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "input_datetime",
		Service: "set_datetime",
		ServiceData: map[string]any{
			"timestamp": fmt.Sprint(value.Unix()),
		},
		Target: message.Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputDatetime) Reload(ctx context.Context) (any, error) {
	req := message.CallServiceData{
		Domain:  "input_datetime",
		Service: "reload",
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
