package services

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

type BaseServiceRequest struct {
	Id          int64          `json:"id"`
	RequestType string         `json:"type"` // must be set to "call_service"
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
