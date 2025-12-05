package services

import (
	"saml.dev/gome-assistant/internal"
)

// API is the interface that the individual services use to interact
// with HomeAssistant.
type API interface {
	WriteMessage(msg any) error
}

func BuildService[
	T AdaptiveLighting |
		AlarmControlPanel |
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
		Timer |
		Vacuum |
		ZWaveJS,
](api API) *T {
	return &T{api: api}
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

func NewBaseServiceRequest(entityId string) BaseServiceRequest {
	id := internal.GetId()
	bsr := BaseServiceRequest{
		Id:          id,
		RequestType: "call_service",
	}
	if entityId != "" {
		bsr.Target.EntityId = entityId
	}
	return bsr
}
