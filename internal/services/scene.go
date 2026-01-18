package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type Scene struct {
	api API
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(
	ctx context.Context, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "scene",
		Service:     "apply",
		ServiceData: serviceData,
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Create a scene entity. Takes an entityID and an optional
// map that is translated into service_data.
func (s Scene) Create(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "scene",
		Service:     "create",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Reload the scenes.
func (s Scene) Reload(ctx context.Context) (any, error) {
	req := message.CallServiceData{
		Domain:  "scene",
		Service: "reload",
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// TurnOn a scene entity. Takes an entityID and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "scene",
		Service:     "turn_on",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
