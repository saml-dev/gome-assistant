package message

import "time"

type StateChangedEventMessage EventMessage[StateChangedData]

type StateChangedData struct {
	EntityID string            `json:"entity_id"`
	NewState StateChangedState `json:"new_state"`
	OldState StateChangedState `json:"old_state"`
}

type StateChangedState struct {
	EntityID    string         `json:"entity_id"`
	LastChanged time.Time      `json:"last_changed"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastUpdated time.Time      `json:"last_updated"`
	Context     EventContext   `json:"context"`
}
