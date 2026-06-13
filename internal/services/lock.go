package services

import "context"

/* Structs */

type Lock struct {
	api API
}

/* Public API */

// Lock a lock entity. Takes entity IDs and an optional service_data,
// which must be serializable to a JSON object.
func (l Lock) Lock(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "lock",
		Service:     "lock",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Unlock a lock entity. Takes entity IDs and an optional
// service_data, which must be serializable to a JSON object.
func (l Lock) Unlock(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "lock",
		Service:     "unlock",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
