package services

import "context"

/* Structs */

type InputNumber struct {
	api API
}

/* Public API */

func (ib InputNumber) Set(
	ctx context.Context, entityID string, value float32,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "input_number",
		Service:     "set_value",
		ServiceData: map[string]any{"value": value},
		Target:      Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputNumber) Increment(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "increment",
		Target:  Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputNumber) Decrement(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "decrement",
		Target:  Entity(entityID),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputNumber) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "reload",
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
