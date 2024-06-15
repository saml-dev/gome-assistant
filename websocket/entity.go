package websocket

// "state_changed" events are compressed in a rather awkward way.
// These types help pick them apart.

type Entity[AttributesT any] struct {
	State       EntityState `json:"state"`
	Attributes  AttributesT `json:"attributes"`
	Context     Context     `json:"context"`
	LastChanged TimeStamp   `json:"last_changed"`
}

type EntityItem[AttributesT any] struct {
	EntityID string `json:"entity_id"`
	Entity[AttributesT]
}

// CompressedEntity is similar to `Entity` except that the JSON field
// names are abbreviated.
type CompressedEntity[AttributesT any] struct {
	State       EntityState `json:"s"`
	Attributes  AttributesT `json:"a"`
	Context     Context     `json:"c"`
	LastChanged TimeStamp   `json:"lc"`
}

// EntityState is the state of an entity ( // E.g., "on", "off",
// "unavailable"; there are probably more.
type EntityState string

func (s EntityState) On() bool {
	return s == "on"
}

func (s EntityState) Off() bool {
	return s == "off"
}

func (s EntityState) Unavailable() bool {
	return s == "unavailable"
}
