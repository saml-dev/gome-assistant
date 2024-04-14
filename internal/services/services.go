package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

func BuildService[
	T AlarmControlPanel |
		Climate |
		Cover |
		Light |
		HomeAssistant |
		Lock |
		MediaPlayer |
		Switch |
		InputBoolean |
		InputButton |
		InputDatetime |
		InputText |
		InputNumber |
		Event |
		Notify |
		Number |
		Scene |
		Script |
		TTS |
		Vacuum |
		ZWaveJS,
](conn *websocket.Conn) *T {
	return &T{conn: conn}
}

type BaseServiceRequest struct {
	Id          int64          `json:"id"`
	RequestType string         `json:"type"` // hardcoded "call_service"
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      struct {
		EntityId string `json:"entity_id,omitempty"`
	} `json:"target,omitempty"`
}

func NewBaseServiceRequest(conn *websocket.Conn, entityId string) BaseServiceRequest {
	id := conn.NextID()
	bsr := BaseServiceRequest{
		Id:          id,
		RequestType: "call_service",
	}
	if entityId != "" {
		bsr.Target.EntityId = entityId
	}
	return bsr
}
