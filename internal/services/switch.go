package services

import "context"

/* Structs */

type Switch struct {
	api API
}

/* Public API */

func (s Switch) TurnOn(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "switch",
		Service: "turn_on",
		Target:  Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s Switch) Toggle(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "switch",
		Service: "toggle",
		Target:  Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s Switch) TurnOff(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "switch",
		Service: "turn_off",
		Target:  Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
