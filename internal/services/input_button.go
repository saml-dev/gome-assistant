package services

import (
	"context"

	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type InputButton struct {
	service Service
}

func NewInputButton(service Service) *InputButton {
	return &InputButton{
		service: service,
	}
}

/* Public API */

func (ib InputButton) Press(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_button", "press",
		nil, EntityTarget(entityID),
	)
}

func (ib InputButton) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_button", "reload", nil, Target{},
	)
}
