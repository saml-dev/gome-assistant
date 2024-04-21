package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
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
func (ha *HomeAssistant) TurnOn(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return ha.service.CallService(
		ctx, "homeassistant", "turn_on",
		serviceData, EntityTarget(entityID),
	)
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return ha.service.CallService(
		ctx, "homeassistant", "toggle",
		serviceData, EntityTarget(entityID),
	)
}

func (ha *HomeAssistant) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ha.service.CallService(
		ctx, "homeassistant", "turn_off",
		nil, EntityTarget(entityID),
	)
}
