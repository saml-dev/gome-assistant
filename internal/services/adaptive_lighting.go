package services

import "context"

/* Structs */

type AdaptiveLighting struct {
	api API
}

/* Public API */

// Set manual control for an adaptive lighting entity.
func (al AdaptiveLighting) SetManualControl(
	ctx context.Context, entityID string, enabled bool,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "adaptive_lighting",
		Service: "set_manual_control",
		ServiceData: map[string]any{
			"entity_id":      entityID,
			"manual_control": enabled,
		},
		Target: Entity(entityID),
	}

	var result any
	if err := al.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
