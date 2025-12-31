package services

import "context"

type Number struct {
	api API
}

func (ib Number) SetValue(
	ctx context.Context, entityID string, value float32,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "number",
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

func (ib Number) MustSetValue(
	ctx context.Context, entityID string, value float32,
) {
	if _, err := ib.SetValue(ctx, entityID, value); err != nil {
		panic(err)
	}
}
