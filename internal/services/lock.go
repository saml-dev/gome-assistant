package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
func (l Lock) Lock(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "lock", "lock",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Unlock a lock entity.
func (l Lock) Unlock(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "lock", "unlock",
		serviceData, ga.EntityTarget(entityID),
	)
}
