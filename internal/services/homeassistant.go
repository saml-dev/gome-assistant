package services

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/http"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
	"nhooyr.io/websocket"
)

type HomeAssistant struct {
	conn       *websocket.Conn
	ctx        context.Context
	httpClient *http.HttpClient
}

func (ha *HomeAssistant) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_on"

	ws.WriteMessage(req, ha.conn, ha.ctx)
}

func (ha *HomeAssistant) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_off"

	ws.WriteMessage(req, ha.conn, ha.ctx)
}

func (ha *HomeAssistant) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "toggle"

	ws.WriteMessage(req, ha.conn, ha.ctx)
}
