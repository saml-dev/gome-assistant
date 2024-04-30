package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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

func (c Climate) SetFanMode(
	target ga.Target, fanMode string,
) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "climate", "set_fan_mode",
		map[string]any{"fan_mode": fanMode},
		target,
	)
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
) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "climate", "set_temperature",
		setTemperatureRequest.ToJSON(),
		target,
	)
}
