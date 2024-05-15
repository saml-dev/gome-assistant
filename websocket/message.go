package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// RawMessage is like `json.RawMessage`, except that its `String()`
// method converts it directly to a string.
type RawMessage json.RawMessage

func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON delegates to `json.RawMessage`. (The method has a
// pointer receiver, so we have to implement it explicitly.)
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	return (*json.RawMessage)(m).UnmarshalJSON(data)
}

func (rm RawMessage) String() string {
	return string(rm)
}

// BaseMessage implements the required part of any websocket message.
// The idea is to embed this type in other message types.
type BaseMessage struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
}

type Request interface {
	GetID() int64
	SetID(id int64)
}

func (msg *BaseMessage) GetID() int64 {
	return msg.ID
}

func (msg *BaseMessage) SetID(id int64) {
	msg.ID = id
}

// Message holds a complete message, only partly parsed. The entire,
// original, unparsed message is available in the `Raw` field.
type Message struct {
	BaseMessage

	// Raw contains the original, full, unparsed message (including
	// fields `Type` and `ID`, which also appear in `BaseMessage`).
	Raw RawMessage `json:"-"`
}

type BaseResultMessage struct {
	BaseMessage
	Success bool         `json:"success"`
	Error   *ResultError `json:"error,omitempty"`
}

type ResultError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *ResultError) Error() string {
	switch {
	case err.Code != "" && err.Message != "":
		return fmt.Sprintf("%s: %s", err.Code, err.Message)
	case err.Code == "" && err.Message != "":
		return fmt.Sprintf("unknown_error: %s", err.Message)
	case err.Code != "" && err.Message == "":
		return fmt.Sprintf("%s", err.Code)
	default:
		// This seems not to be an error at all.
		return fmt.Sprintf("INVALID (seems not to be an error)")
	}
}

type ResultMessage struct {
	BaseResultMessage

	// Raw contains the "result" part of the message, unparsed.
	Result RawMessage `json:"result"`
}

// GetResult parses a result out of `msg` (which must have type
// "result"). If `msg` indicates that an error occurred, convert that
// to an error and return it. Parse the result into `result`, which
// must be unmarshalable as JSON.
func (msg Message) GetResult(result any) error {
	if msg.Type != "result" {
		return fmt.Errorf(
			"response message was not of type 'result': %#v", msg,
		)
	}
	var resultMsg ResultMessage
	if err := json.Unmarshal(msg.Raw, &resultMsg); err != nil {
		return fmt.Errorf("unmarshaling result message: %w", err)
	}
	if !resultMsg.Success {
		if resultMsg.Error == nil {
			return fmt.Errorf(
				"request did not succeed but no error was returned",
			)
		}
		return resultMsg.Error
	}

	if err := json.Unmarshal(resultMsg.Result, result); err != nil {
		return fmt.Errorf("unmarshalling result from %q: %w", resultMsg.Result, err)
	}
	return nil
}

// "state_changed" events are compressed in a rather awkward way.
// These types help pick them apart.

type EntityState struct {
	State       RawMessage            `json:"state"`
	Attributes  map[string]RawMessage `json:"attributes"`
	Context     RawMessage            `json:"context"`
	LastChanged RawMessage            `json:"last_changed"`
}

type EntityStateItem struct {
	EntityID string `json:"entity_id"`
	EntityState
}

// CompressedEntityState is similar to `EntityState` except that it
// doesn't include the entity ID and the JSON field names are
// abbreviated.
type CompressedEntityState struct {
	State       RawMessage            `json:"s"`
	Attributes  map[string]RawMessage `json:"a"`
	Context     RawMessage            `json:"c"`
	LastChanged RawMessage            `json:"lc"`
}

// CompressedEntityChange keeps tracks of fields added and removed as
// part of a change. Fields that are mutated appear as "additions".
type CompressedEntityChange struct {
	Additions CompressedEntityState `json:"+,omitempty"`
	Removals  struct {
		Attributes []string `json:"a"`
		Context    []string `json:"c"`
	} `json:"-,omitempty"`
}

type CompressedStateChangedMessage struct {
	BaseMessage
	Event struct {
		Added   map[string]CompressedEntityState  `json:"a,omitempty"`
		Changed map[string]CompressedEntityChange `json:"c,omitempty"`
		Removed []string                          `json:"r,omitempty"`
	} `json:"event"`
}

// Apply applies the changes indicated in `msg` to the entity with the
// specified `entityID` whose old state was `oldState`, returning the
// new state.
func (msg CompressedStateChangedMessage) Apply(
	entityID string, oldState EntityState,
) (EntityState, error) {
	if state, ok := msg.Event.Added[entityID]; ok {
		// This entityID was added. The new state was right there in
		// the message.
		return EntityState(state), nil
	}
	if change, ok := msg.Event.Changed[entityID]; ok {
		state := oldState.State
		if len(change.Additions.State) != 0 {
			state = change.Additions.State
		}
		// The existing entry has had some fields changed.
		return EntityState{
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
			return EntityState{}, nil
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
