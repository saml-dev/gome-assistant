package websocket

import (
	"time"
)

// "state_changed" events are compressed in a rather awkward way.
// These types help pick them apart.

type Entity struct {
	State       EntityState           `json:"state"`
	Attributes  map[string]RawMessage `json:"attributes"`
	Context     RawMessage            `json:"context"`
	LastChanged time.Time             `json:"last_changed"`
}

type EntityItem struct {
	EntityID string `json:"entity_id"`
	Entity
}

// CompressedEntity is similar to `Entity` except that the JSON field
// names are abbreviated.
type CompressedEntity struct {
	State       EntityState           `json:"s"`
	Attributes  map[string]RawMessage `json:"a"`
	Context     RawMessage            `json:"c"`
	LastChanged time.Time             `json:"lc"`
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
