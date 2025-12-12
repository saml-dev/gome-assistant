package services

/* Structs */

type AdaptiveLighting struct {
	api API
}

/* Public API */

// Set manual control for an adaptive lighting entity.
func (al AdaptiveLighting) SetManualControl(entityID string, enabled bool) error {
	req := BaseServiceRequest{
		Domain:  "adaptive_lighting",
		Service: "set_manual_control",
		ServiceData: map[string]any{
			"entity_id":      entityID,
			"manual_control": enabled,
		},
		Target: Entity(entityID),
	}

	return al.api.Call(req)
}
