package message

// CallServiceRequest represents a message that can be sent to request
// an API call. Its `Type` field must be set to "call_service".
type CallServiceRequest struct {
	BaseMessage
	CallServiceData
}

// CallServiceData contains the fields needed to make an HA API call.
// `Target` is not always required. The `ServiceData` field can
// contain arbitrary data needed for a particular call.
type CallServiceData struct {
	Domain      string `json:"domain"`
	Service     string `json:"service"`
	Target      Target `json:"target,omitempty"`
	ServiceData any    `json:"service_data,omitempty"`
}
