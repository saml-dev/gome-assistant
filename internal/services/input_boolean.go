package services

import "context"

/* Structs */

type InputBoolean struct {
	api API
}

/* Public API */

func (ib InputBoolean) TurnOn(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_on",
		Target:  Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputBoolean) Toggle(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "toggle",
		Target:  Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputBoolean) TurnOff(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_off",
		Target:  Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputBoolean) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "reload",
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
