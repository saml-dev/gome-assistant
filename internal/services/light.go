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
	ctx context.Context, entityIDs []string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entities(entityIDs),
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
	ctx context.Context, entityIDs []string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "toggle",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entities(entityIDs),
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (l Light) TurnOff(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}
	return result, nil
}
