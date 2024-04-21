package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
)

// Target represents the target of the service call, if applicable.
type Target struct {
	EntityID string `json:"entity_id,omitempty"`
}

func EntityTarget(entityID string) Target {
	return Target{
		EntityID: entityID,
	}
}

type Service interface {
	Call(
		ctx context.Context, req websocket.Request,
	) (websocket.Message, error)

	CallService(
		ctx context.Context, domain string, service string, serviceData any, target Target,
	) (websocket.Message, error)
}
