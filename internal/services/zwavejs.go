package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type ZWaveJS struct {
	api API
}

/* Public API */

// ZWaveJS bulk_set_partial_config_parameters service.
func (zw ZWaveJS) BulkSetPartialConfigParam(
	ctx context.Context, entityID string, parameter int, value any,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "zwave_js",
		Service: "bulk_set_partial_config_parameters",
		ServiceData: map[string]any{
			"parameter": parameter,
			"value":     value,
		},
		Target: message.Entity(entityID),
	}

	var result any
	if err := zw.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
