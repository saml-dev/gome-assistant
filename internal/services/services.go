package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

type Service interface {
	Call(
		ctx context.Context, req websocket.Request, result any,
	) error

	CallService(
		ctx context.Context, domain string, service string, serviceData any, target ga.Target,
		result any,
	) error
}
