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

// CallService is a type that always serializes as `"call_service"`.
type CallService struct{}

func (CallService) String() string {
	return "call_service"
}

func (CallService) MarshalJSON() ([]byte, error) {
	return []byte(`"call_service"`), nil
}

type CallServiceRequest struct {
	ID          int64       `json:"id"`
	RequestType CallService `json:"type"` // hardcoded "call_service"
	Domain      string      `json:"domain"`
	Service     string      `json:"service"`

	// ServiceData must be serializable to a JSON object.
	ServiceData any `json:"service_data,omitempty"`

	Target Target `json:"target,omitempty"`
}
