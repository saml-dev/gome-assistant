package services

import "context"

/* Structs */

type ZWaveJS struct {
	api API
}

/* Public API */

// ZWaveJS bulk_set_partial_config_parameters service.
func (zw ZWaveJS) BulkSetPartialConfigParam(
	ctx context.Context, target Target, parameter int, value any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "zwave_js",
		Service: "bulk_set_partial_config_parameters",
		ServiceData: map[string]any{
			"parameter": parameter,
			"value":     value,
		},
		Target: target,
	}

	var result any
	if err := zw.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
