package services

import "context"

/* Structs */

type Light struct {
	api API
}

/* Public API */

// TurnOn a light entity. Takes entity IDs and an optional
// service_data, which must be serializable to a JSON object.
func (l Light) TurnOn(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a light entity. Takes entity IDs and an optional
// service_data, which must be serializable to a JSON object.
func (l Light) Toggle(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "toggle",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (l Light) TurnOff(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target:  target,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}
	return result, nil
}
