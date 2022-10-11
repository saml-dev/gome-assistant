package gomeassistant

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/http"
	"github.com/saml-dev/gome-assistant/internal/services"
	"nhooyr.io/websocket"
)

type Service struct {
	HomeAssistant *services.HomeAssistant
	Light         *services.Light
}

func NewService(conn *websocket.Conn, ctx context.Context, httpClient *http.HttpClient) *Service {
	return &Service{
		Light:         services.BuildService[services.Light](conn, ctx),
		HomeAssistant: services.BuildService[services.HomeAssistant](conn, ctx),
	}
}
