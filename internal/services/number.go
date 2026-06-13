package services

import "context"

type Number struct {
	api API
}

func (ib Number) SetValue(
	ctx context.Context, entityIDs []string, value float32,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "number",
		Service:     "set_value",
		ServiceData: map[string]any{"value": value},
		Target:      Entities(entityIDs),
	}

	var result any
	if err := ib.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ib Number) MustSetValue(
	ctx context.Context, entityIDs []string, value float32,
) {
	if _, err := ib.SetValue(ctx, entityIDs, value); err != nil {
		panic(err)
	}
}
