package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type InputText struct {
	service Service
}

func NewInputText(service Service) *InputText {
	return &InputText{
		service: service,
	}
}

/* Public API */

func (ib InputText) Set(target ga.Target, value string) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_text", "set_value",
		map[string]any{
			"value": value,
		},
		target,
	)
}

func (ib InputText) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_text", "reload", nil, ga.Target{},
	)
}
