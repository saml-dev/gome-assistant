package services

import "context"

/* Structs */

type Light struct {
	api API
}

/* Public API */

// TurnOn a light entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Light) TurnOn(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "turn_on",
		Target:      Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a light entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Light) Toggle(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "toggle",
		Target:      Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (l Light) TurnOff(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target:  Entity(entityID),
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}
	return result, nil
}
