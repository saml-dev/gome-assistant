package services

import "saml.dev/gome-assistant/websocket"

// API is the interface that the individual services use to interact
// with HomeAssistant.
type API interface {
	Call(req BaseServiceRequest) error
	FireEvent(eventType string, eventData map[string]any) error
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

// CallServiceMessage represents a message that can be sent to request
// an API call. Its `Type` field must be set to "call_service".
type CallServiceMessage struct {
	websocket.BaseMessage
	BaseServiceRequest
}

// BaseServiceRequest contains the fields needed to make an HA API
// call. `ServiceData` can contain arbitrary data needed for a
// particular call.
type BaseServiceRequest struct {
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      Target         `json:"target,omitempty"`
}

type Target struct {
	EntityID string `json:"entity_id,omitempty"`
}

func Entity(entityID string) Target {
	return Target{
		EntityID: entityID,
	}
}
