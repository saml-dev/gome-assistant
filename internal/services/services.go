package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
)

type Service interface {
	Call(
		ctx context.Context, req websocket.Request,
	) (websocket.Message, error)

	CallService(
		ctx context.Context, domain string, service string, serviceData any, target ga.Target,
	) (websocket.Message, error)
}
