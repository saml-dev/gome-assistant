package services

// CallService is a type that always serializes as `"call_service"`.
type CallService struct{}

func (CallService) String() string {
	return "call_service"
}

func (CallService) MarshalJSON() ([]byte, error) {
	return []byte(`"call_service"`), nil
}

// Target represents the target of the service call, if applicable.
type Target struct {
	EntityID string `json:"entity_id,omitempty"`
}

type CallServiceRequest struct {
	ID          int64          `json:"id"`
	RequestType CallService    `json:"type"` // hardcoded "call_service"
	Domain      string         `json:"domain"`
	Service     string         `json:"service"`
	ServiceData map[string]any `json:"service_data,omitempty"`
	Target      Target         `json:"target,omitempty"`
}
