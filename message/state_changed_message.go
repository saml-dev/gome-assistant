package message

import "time"

type StateChangedMessage struct {
	BaseMessage
	Event struct {
		Data      StateData `json:"data"`
		EventType string    `json:"event_type"`
		Origin    string    `json:"origin"`
	} `json:"event"`
}

type StateData struct {
	EntityID string       `json:"entity_id"`
	NewState MessageState `json:"new_state"`
	OldState MessageState `json:"old_state"`
}

type MessageState struct {
	EntityID    string         `json:"entity_id"`
	LastChanged time.Time      `json:"last_changed"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
}
