package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type Lock struct {
	api API
}

/* Public API */

// Lock a lock entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Lock) Lock(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "lock",
		Service:     "lock",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Unlock a lock entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Lock) Unlock(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "lock",
		Service:     "unlock",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := l.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
