package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type HomeAssistant struct {
	service Service
}

func NewHomeAssistant(service Service) *HomeAssistant {
	return &HomeAssistant{
		service: service,
	}
}

// TurnOn a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := ha.service.CallService(
		ctx, "homeassistant", "turn_on",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := ha.service.CallService(
		ctx, "homeassistant", "toggle",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ha *HomeAssistant) TurnOff(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := ha.service.CallService(
		ctx, "homeassistant", "turn_off",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
