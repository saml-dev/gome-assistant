package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type Number struct {
	service Service
}

func NewNumber(service Service) *Number {
	return &Number{
		service: service,
	}
}

/* Public API */

func (ib Number) SetValue(entityID string, value float32) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "number", "set_value",
		map[string]any{"value": value},
		ga.EntityTarget(entityID),
	)
}
