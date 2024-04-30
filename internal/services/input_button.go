package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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

func (ib InputButton) Press(target ga.Target) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_button", "press",
		nil, target,
	)
}

func (ib InputButton) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_button", "reload", nil, ga.Target{},
	)
}
