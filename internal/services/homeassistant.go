package services

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/http"
	"nhooyr.io/websocket"
)

type HomeAssistant struct {
	conn       *websocket.Conn
	ctx        context.Context
	httpClient *http.HttpClient
}

// TODO: design how much reuse I can get between request types. E.g.
// only difference between light.turnon and homeassistant.turnon is
// domain and extra data
