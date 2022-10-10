package services

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/setup"
	"nhooyr.io/websocket"
)

type Light struct {
	conn *websocket.Conn
	ctx  context.Context
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

func LightOnRequest(entityId string) LightRequest {
	req := LightRequest{
		Id:      5,
		Type:    "call_service",
		Domain:  "light",
		Service: "turn_on",
	}
	req.Target.EntityId = entityId
	return req
}

func LightOffRequest(entityId string) LightRequest {
	req := LightOnRequest(entityId)
	req.Service = "turn_off"
	return req
}

func (l Light) TurnOn(entityId string) {
	req := LightOnRequest(entityId)
	setup.WriteMessage(req, l.conn, l.ctx)
}

func (l Light) TurnOff(entityId string) {
	req := LightOffRequest(entityId)
	setup.WriteMessage(req, l.conn, l.ctx)
}
