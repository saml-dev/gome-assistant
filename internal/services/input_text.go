package services

import "context"

/* Structs */

type InputText struct {
	api API
}

/* Public API */

func (ib InputText) Set(
	ctx context.Context, target Target, value string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_text",
		Service: "set_value",
		ServiceData: map[string]any{
			"value": value,
		},
		Target: target,
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib InputText) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "input_text",
		Service: "reload",
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
