package services

import "context"

/* Structs */

type InputButton struct {
	api API
}

/* Public API */

func (ib InputButton) Press(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_button",
		Service: "press",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputButton) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_button",
		Service: "reload",
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
