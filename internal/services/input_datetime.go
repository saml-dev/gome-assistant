package services

import (
	"context"
	"fmt"
	"time"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputDatetime struct {
	service Service
}

func NewInputDatetime(service Service) *InputDatetime {
	return &InputDatetime{
		service: service,
	}
}

/* Public API */

func (ib InputDatetime) Set(entityID string, value time.Time) (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_datetime", "set_datetime",
		map[string]any{
			"timestamp": fmt.Sprint(value.Unix()),
		},
		EntityTarget(entityID),
	)
}

func (ib InputDatetime) Reload() (websocket.Message, error) {
	ctx := context.TODO()
	return ib.service.CallService(
		ctx, "input_datetime", "reload", nil, Target{},
	)
}
