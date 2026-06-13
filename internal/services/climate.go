package services

import (
	"context"

	"saml.dev/gome-assistant/types"
)

/* Structs */

type Climate struct {
	api API
}

/* Public API */

func (c Climate) SetFanMode(
	ctx context.Context, entityIDs []string, fanMode string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "climate",
		Service:     "set_fan_mode",
		ServiceData: map[string]any{"fan_mode": fanMode},
		Target:      Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c Climate) SetTemperature(
	ctx context.Context, entityIDs []string, serviceData types.SetTemperatureRequest,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "climate",
		Service:     "set_temperature",
		ServiceData: serviceData.ToJSON(),
		Target:      Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
