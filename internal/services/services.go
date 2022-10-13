package services

import (
	"context"
	"fmt"

	"nhooyr.io/websocket"
)

func BuildService[
	T Light |
		HomeAssistant |
		Lock |
		Switch |
		InputBoolean |
		InputButton |
		InputDatetime |
		InputText |
		InputNumber,
](conn *websocket.Conn, ctx context.Context) *T {
	return &T{conn: conn, ctx: ctx}
}

type BaseServiceRequest struct {
	Id          string         `json:"id"`
	RequestType string         `json:"type"` // hardcoded "call_service"
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      struct {
		EntityId string `json:"entity_id"`
	} `json:"target,omitempty"`
}

var id int64 = 1

func NewBaseServiceRequest(entityId string) BaseServiceRequest {
	bsr := BaseServiceRequest{
		Id:          fmt.Sprint(id),
		RequestType: "call_service",
	}
	if entityId != "" {
		bsr.Target.EntityId = entityId
	}
	id += 1
	return bsr
}
