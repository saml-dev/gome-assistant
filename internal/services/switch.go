package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type Switch struct {
	api API
}

/* Public API */

func (s Switch) TurnOn(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "switch",
		Service: "turn_on",
		Target:  message.Entity(entityID),
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
	req := message.CallServiceData{
		Domain:  "switch",
		Service: "toggle",
		Target:  message.Entity(entityID),
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
	req := message.CallServiceData{
		Domain:  "switch",
		Service: "turn_off",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := s.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
