package services

/* Structs */

type ZWaveJS struct {
	api API
}

/* Public API */

// ZWaveJS bulk_set_partial_config_parameters service.
func (zw ZWaveJS) BulkSetPartialConfigParam(entityID string, parameter int, value any) error {
	req := BaseServiceRequest{
		Domain:  "zwave_js",
		Service: "bulk_set_partial_config_parameters",
		ServiceData: map[string]any{
			"parameter": parameter,
			"value":     value,
		},
		Target: Entity(entityID),
	}
	return zw.api.CallAndForget(req)
}
