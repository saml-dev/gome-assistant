package services

/* Structs */

type AdaptiveLighting struct {
	api API
}

/* Public API */

// Set manual control for an adaptive lighting entity.
func (al AdaptiveLighting) SetManualControl(entityId string, enabled bool) error {
	req := BaseServiceRequest{
		Domain:  "adaptive_lighting",
		Service: "set_manual_control",
		ServiceData: map[string]any{
			"entity_id":      entityId,
			"manual_control": enabled,
		},
		Target: Entity(entityId),
	}

	return al.api.Call(req)
}
