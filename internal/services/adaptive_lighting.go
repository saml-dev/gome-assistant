package services

/* Structs */

type AdaptiveLighting struct {
	api API
}

/* Public API */

// Set manual control for an adaptive lighting entity.
func (al AdaptiveLighting) SetManualControl(entityId string, enabled bool) error {
	req := NewBaseServiceRequest("")
	req.Domain = "adaptive_lighting"
	req.Service = "set_manual_control"
	req.ServiceData = map[string]any{
		"entity_id":      entityId,
		"manual_control": enabled,
	}

	return al.api.WriteMessage(req)
}
