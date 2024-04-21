package services

import (
	"context"

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
		serviceData, EntityTarget(entityID),
	)
}

// Unlock a lock entity.
func (l Lock) Unlock(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return l.service.CallService(
		ctx, "lock", "unlock",
		serviceData, EntityTarget(entityID),
	)
}
