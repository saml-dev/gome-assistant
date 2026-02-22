package message

// EntityState is the state of an entity in Home Assistant, reported
// by the `/api/states` and `/api/states/<entity_id>` REST APIs.
type EntityState struct {
	EntityID    string         `json:"entity_id"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastChanged TimeStamp      `json:"last_changed"`
	LastUpdated TimeStamp      `json:"last_updated"`
}
