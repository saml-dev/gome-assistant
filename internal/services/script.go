package services

import "context"

/* Structs */

type Script struct {
	api API
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "reload",
		Target:  Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "toggle",
		Target:  Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// TurnOff a script that was created in the HA UI.
func (s Script) TurnOff(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "turn_off",
		Target:  Entity(""),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// TurnOn a script that was created in the HA UI.
func (s Script) TurnOn(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "turn_on",
		Target:  Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
