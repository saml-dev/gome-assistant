package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type Script struct {
	service Service
}

func NewScript(service Service) *Script {
	return &Script{
		service: service,
	}
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "script", "reload",
		nil, ga.EntityTarget(entityID),
	)
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "script", "toggle",
		nil, ga.EntityTarget(entityID),
	)
}

// Turn off a script that was created in the HA UI.
func (s Script) TurnOff() (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "script", "turn_off",
		nil, ga.Target{},
	)
}

// Turn on a script that was created in the HA UI.
func (s Script) TurnOn(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return s.service.CallService(
		ctx, "script", "turn_on",
		nil, ga.EntityTarget(entityID),
	)
}
