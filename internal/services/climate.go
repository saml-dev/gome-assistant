package services

import (
	"saml.dev/gome-assistant/types"
)

/* Structs */

type Climate struct {
	api API
}

/* Public API */

func (c Climate) SetFanMode(entityID string, fanMode string) error {
	req := BaseServiceRequest{
		Domain:      "climate",
		Service:     "set_fan_mode",
		ServiceData: map[string]any{"fan_mode": fanMode},
		Target:      Entity(entityID),
	}
	return c.api.Call(req)
}

func (c Climate) SetTemperature(entityID string, serviceData types.SetTemperatureRequest) error {
	req := BaseServiceRequest{
		Domain:      "climate",
		Service:     "set_temperature",
		ServiceData: serviceData.ToJSON(),
		Target:      Entity(entityID),
	}
	return c.api.Call(req)
}
