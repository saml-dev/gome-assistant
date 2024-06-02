package websocket

import (
	"encoding/json"
	"log/slog"
	"time"
)

// "state_changed" events are compressed in a rather awkward way.
// These types help pick them apart.

type Entity struct {
	State       RawMessage            `json:"state"`
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
	State       RawMessage            `json:"s"`
	Attributes  map[string]RawMessage `json:"a"`
	Context     RawMessage            `json:"c"`
	LastChanged time.Time             `json:"lc"`
}

// CompressedEntityChange keeps tracks of fields added and removed as
// part of a change. Fields that are mutated appear as "additions".
type CompressedEntityChange struct {
	Additions CompressedEntity `json:"+,omitempty"`
	Removals  struct {
		Attributes []string `json:"a"`
		Context    []string `json:"c"`
	} `json:"-,omitempty"`
}

type CompressedStateChangedMessage struct {
	BaseMessage
	Event struct {
		Added   map[string]CompressedEntity       `json:"a,omitempty"`
		Changed map[string]CompressedEntityChange `json:"c,omitempty"`
		Removed []string                          `json:"r,omitempty"`
	} `json:"event"`
}

// Apply applies the changes indicated in `msg` to the entity with the
// specified `entityID` whose old state was `oldState`, returning the
// new state.
func (msg CompressedStateChangedMessage) Apply(
	entityID string, oldState Entity,
) (Entity, error) {
	if state, ok := msg.Event.Added[entityID]; ok {
		// This entityID was added. The new state was right there in
		// the message.
		return Entity(state), nil
	}
	if change, ok := msg.Event.Changed[entityID]; ok {
		state := oldState.State
		if len(change.Additions.State) != 0 {
			state = change.Additions.State
		}
		// The existing entry has had some fields changed.
		return Entity{
			State: state,
			Attributes: mergeMaps(
				oldState.Attributes,
				change.Additions.Attributes,
				change.Removals.Attributes,
			),
			// FIXME: apparently, context can also be a single string.
			Context: mergeContexts(
				oldState.Context,
				change.Additions.Context,
				change.Removals.Context,
			),
			LastChanged: change.Additions.LastChanged,
		}, nil
	}
	for _, eid := range msg.Event.Removed {
		if eid == entityID {
			return Entity{}, nil
		}
	}
	return oldState, nil
}

func mergeMaps(old, additions map[string]RawMessage, removals []string) map[string]RawMessage {
	new := make(map[string]RawMessage, len(old)+len(additions)-len(removals))
	for k, v := range old {
		new[k] = v
	}
	for k, v := range additions {
		new[k] = v
	}
	for _, k := range removals {
		delete(new, k)
	}
	return new
}

func mergeContexts(old, additions RawMessage, removals []string) RawMessage {
	switch {
	case len(old) == 0:
		return additions
	case old[0] == '"':
		// The context is a single string.
		if len(additions) != 0 {
			return additions
		}
		return old
	case old[0] == '{':
		// The context is an object.
		var contextMap map[string]RawMessage
		if err := json.Unmarshal(old, &contextMap); err != nil {
			slog.Error("cannot unmarshal old context",
				"old", string(old),
				"error", err,
			)
			return additions
		}
		var addMap map[string]RawMessage
		if len(additions) != 0 {
			if err := json.Unmarshal(additions, &addMap); err != nil {
				slog.Error("cannot unmarshal additions",
					"additions", string(additions),
					"error", err,
				)
				return old
			}
		}
		newMap := mergeMaps(contextMap, addMap, removals)
		newContext, err := json.Marshal(newMap)
		if err != nil {
			slog.Error("cannot marshal new context",
				"context", newMap,
				"error", err,
			)
			return old
		}
		return newContext
	default:
		return old
	}
}
