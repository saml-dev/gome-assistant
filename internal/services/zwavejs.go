package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type ZWaveJS struct {
	service Service
}

func NewZWaveJS(service Service) *ZWaveJS {
	return &ZWaveJS{
		service: service,
	}
}

/* Public API */

// ZWaveJS bulk_set_partial_config_parameters service.
func (zw ZWaveJS) BulkSetPartialConfigParam(
	entityID string, parameter int, value any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return zw.service.CallService(
		ctx, "zwave_js", "bulk_set_partial_config_parameters",
		map[string]any{
			"parameter": parameter,
			"value":     value,
		},
		EntityTarget(entityID),
	)
}
