package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Lock struct {
	service Service
}

func NewLock(service Service) *Lock {
	return &Lock{
		service: service,
	}
}

/* Public API */

// Lock a lock entity.
func (l Lock) Lock(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := l.service.CallService(
		ctx, "lock", "lock",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Unlock a lock entity.
func (l Lock) Unlock(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := l.service.CallService(
		ctx, "lock", "unlock",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
