package services

import "context"

/* Structs */

type Scene struct {
	api API
}

/* Public API */

// Apply a scene. Takes an optional service_data, which must be
// serializable to a JSON object.
func (s Scene) Apply(
	ctx context.Context, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "apply",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(""),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Create a scene entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (s Scene) Create(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "create",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Reload the scenes.
func (s Scene) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "scene",
		Service: "reload",
		Target:  Entity(""),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// TurnOn a scene entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (s Scene) TurnOn(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
