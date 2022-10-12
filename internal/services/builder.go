package services

import (
	"context"

	"github.com/saml-dev/gome-assistant/internal/http"
	"nhooyr.io/websocket"
)

func BuildService[T Light | HomeAssistant](conn *websocket.Conn, ctx context.Context, httpClient *http.HttpClient) *T {
	return &T{conn: conn, ctx: ctx, httpClient: httpClient}
}

type BaseServiceRequest struct {
	Id          int            `json:"id"`
	RequestType string         `json:"type"` // hardcoded "call_service"
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      struct {
		EntityId string `json:"entity_id"`
	} `json:"target"`
}

func NewBaseServiceRequest(entityId string) BaseServiceRequest {
	bsr := BaseServiceRequest{
		Id:          10,
		RequestType: "call_service",
	}
	bsr.Target.EntityId = entityId
	return bsr
}
