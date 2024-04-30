package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
func (ha *HomeAssistant) TurnOn(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return ha.service.CallService(
		ctx, "homeassistant", "turn_on",
		serviceData, target,
	)
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(target ga.Target, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return ha.service.CallService(
		ctx, "homeassistant", "toggle",
		serviceData, target,
	)
}

func (ha *HomeAssistant) TurnOff(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ha.service.CallService(
		ctx, "homeassistant", "turn_off",
		nil, target,
	)
}
