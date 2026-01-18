package message

type StateChangedEventMessage EventMessage[StateChangedData]

type StateChangedData struct {
	EntityID string            `json:"entity_id"`
	NewState StateChangedState `json:"new_state"`
	OldState StateChangedState `json:"old_state"`
}

type StateChangedState struct {
	EntityID    string         `json:"entity_id"`
	LastChanged TimeStamp      `json:"last_changed"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastUpdated TimeStamp      `json:"last_updated"`
	Context     EventContext   `json:"context"`
}
