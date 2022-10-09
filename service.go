package gomeassistant

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/services"
	"nhooyr.io/websocket"
)

type Service struct {
	HomeAssistant homeAssistant
	Light         services.Light
}

type homeAssistant struct {
	conn websocket.Conn
	ctx  context.Context
}

// type light struct {
// 	conn websocket.Conn
// 	ctx  context.Context
// }
