package services

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/http"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
	"nhooyr.io/websocket"
)

type Light struct {
	conn       *websocket.Conn
	ctx        context.Context
	httpClient *http.HttpClient
}

type LightRequest struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Domain  string `json:"domain"`
	Service string `json:"service"`
	Target  struct {
		EntityId string `json:"entity_id"`
	} `json:"target"`
}

/* Public API */

func (l Light) TurnOn(entityId string) {
	req := newLightOnRequest(entityId)
	ws.WriteMessage(req, l.conn, l.ctx)
}

func (l Light) TurnOff(entityId string) {
	req := newLightOffRequest(entityId)
	ws.WriteMessage(req, l.conn, l.ctx)
}

/* Internal */

func newLightOnRequest(entityId string) LightRequest {
	req := LightRequest{
		Id:      5,
		Type:    "call_service",
		Domain:  "light",
		Service: "turn_on",
	}
	req.Target.EntityId = entityId
	return req
}

func newLightOffRequest(entityId string) LightRequest {
	req := newLightOnRequest(entityId)
	req.Service = "turn_off"
	return req
}
