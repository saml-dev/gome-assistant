package services

import (
	"context"
	"log/slog"

	"saml.dev/gome-assistant/message"
)

// API is the interface that the individual services use to interact
// with HomeAssistant.
type API interface {
	// CallAndForget makes a call to the Home Assistant API but
	// doesn't subscribe to or wait for a response.
	CallAndForget(req BaseServiceRequest) error

	// Call makes a call to the Home Assistant API and waits for a
	// response. The result is unmarshaled into invokes `result`.
	// `result` must be something that `json.Unmarshal()` can
	// deserialize into; typically, it is a pointer. If the result
	// indicates a failure (success==false), then return that as a
	// `*websocket.ResultError`. If another error occurs (e.g.,
	// sending the request or if `ctx` expires), return that error.
	Call(ctx context.Context, req BaseServiceRequest, result any) error

	FireEvent(eventType string, eventData any) error
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
	message.BaseMessage
	BaseServiceRequest
}

// BaseServiceRequest contains the fields needed to make an HA API
// call. `ServiceData` can contain arbitrary data needed for a
// particular call.
type BaseServiceRequest struct {
	Domain      string `json:"domain"`
	Service     string `json:"service"`
	ServiceData any    `json:"service_data,omitempty"`
	Target      Target `json:"target,omitempty"`
}

type Target struct {
	EntityID string `json:"entity_id,omitempty"`
}

func Entity(entityID string) Target {
	return Target{
		EntityID: entityID,
	}
}

func optionalServiceData(serviceData ...any) any {
	switch len(serviceData) {
	case 0:
		return nil
	case 1:
		return serviceData[0]
	default:
		slog.Warn(
			"multiple arguments passed as service data; only the first used",
		)
		return serviceData[0]
	}
}
