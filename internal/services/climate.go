package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Climate struct {
	service Service
}

func NewClimate(service Service) *Climate {
	return &Climate{
		service: service,
	}
}

func (c Climate) SetFanMode(target ga.Target, fanMode string) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "climate", "set_fan_mode",
		map[string]any{"fan_mode": fanMode},
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type SetTemperatureRequest struct {
	Temperature    *float32
	TargetTempHigh *float32
	TargetTempLow  *float32
	HvacMode       string
}

func (r *SetTemperatureRequest) ToJSON() map[string]any {
	m := map[string]any{}
	if r.Temperature != nil {
		m["temperature"] = *r.Temperature
	}
	if r.TargetTempHigh != nil {
		m["target_temp_high"] = *r.TargetTempHigh
	}
	if r.TargetTempLow != nil {
		m["target_temp_low"] = *r.TargetTempLow
	}
	if r.HvacMode != "" {
		m["hvac_mode"] = r.HvacMode
	}
	return m
}

func (c Climate) SetTemperature(
	target ga.Target, setTemperatureRequest SetTemperatureRequest,
) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "climate", "set_temperature",
		setTemperatureRequest.ToJSON(),
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
