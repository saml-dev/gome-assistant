package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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
	target ga.Target, parameter int, value any,
) (any, error) {
	ctx := context.TODO()
	var result any
	err := zw.service.CallService(
		ctx, "zwave_js", "bulk_set_partial_config_parameters",
		map[string]any{
			"parameter": parameter,
			"value":     value,
		},
		target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
