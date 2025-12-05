package services

import (
	"saml.dev/gome-assistant/types"
)

/* Structs */

type Climate struct {
	api API
}

/* Public API */

func (c Climate) SetFanMode(entityId string, fanMode string) error {
	req := BaseServiceRequest{
		Domain:      "climate",
		Service:     "set_fan_mode",
		ServiceData: map[string]any{"fan_mode": fanMode},
		Target:      Entity(entityId),
	}
	return c.api.Call(req)
}

func (c Climate) SetTemperature(entityId string, serviceData types.SetTemperatureRequest) error {
	req := BaseServiceRequest{
		Domain:      "climate",
		Service:     "set_temperature",
		ServiceData: serviceData.ToJSON(),
		Target:      Entity(entityId),
	}
	return c.api.Call(req)
}
