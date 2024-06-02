package websocket

import (
	"encoding/json"
	"fmt"
)

// "state_changed" events are compressed in a rather awkward way.
// These types help pick them apart.

// CompressedEntityChange keeps tracks of fields added and removed as
// part of a change. Fields that are mutated appear as "additions".
type CompressedEntityChange struct {
	Additions CompressedEntity[RawObject] `json:"+,omitempty"`
	Removals  struct {
		Attributes []string `json:"a"`
		Context    []string `json:"c"`
	} `json:"-,omitempty"`
}

type CompressedStateChangedMessage struct {
	BaseMessage
	Event struct {
		Added   map[string]CompressedEntity[RawObject] `json:"a,omitempty"`
		Changed map[string]CompressedEntityChange      `json:"c,omitempty"`
		Removed []string                               `json:"r,omitempty"`
	} `json:"event"`
}

// ApplyChange applies the changes indicated in `msg` to the entity with the
// specified `entityID` and whose old state was `oldEntity`, returning the
// new entity. If the entity was removed altogether, return an empty
// entity.
//
// Because the entity being changed might not store its attributes as
// a generic `RawObject`, we have to do the conversion in an awkward
// way to avoiding needing specialized code for each `AttributeT`:

//  1. Convert the old attributes from an `AttributeT` into a
//     `RawObject`;
//  2. Apply the attribute changes to the `RawObject`;
//  3. Convert the updated `RawObject` back into an `AttributeT`.
func ApplyChange[AttributeT any](
	msg CompressedStateChangedMessage,
	entityID string, oldEntity Entity[AttributeT],
) (Entity[AttributeT], error) {
	for _, eid := range msg.Event.Removed {
		if eid == entityID {
			return Entity[AttributeT]{}, nil
		}
	}

	if entity, ok := msg.Event.Added[entityID]; ok {
		// This entityID was added. The new state was right there in
		// the message.
		var newAttributes AttributeT
		if err := convertTypes(&newAttributes, entity.Attributes); err != nil {
			return Entity[AttributeT]{}, fmt.Errorf(
				"converting the added attributes: %w", err,
			)
		}
		return Entity[AttributeT]{
			State:      entity.State,
			Attributes: newAttributes,
			// FIXME: apparently, context can also be a single string.
			Context:     entity.Context,
			LastChanged: entity.LastChanged,
		}, nil
	}

	change, ok := msg.Event.Changed[entityID]
	if !ok {
		// There were no changes.
		return oldEntity, nil
	}

	// The existing entry has had some fields changed. Apply them to
	// `entity` to produce the new entity:

	newEntity := Entity[AttributeT]{
		State: oldEntity.State,
		Context: mergeContexts(
			oldEntity.Context,
			change.Additions.Context,
			change.Removals.Context,
		),
		LastChanged: change.Additions.LastChanged,
	}

	if change.Additions.State != "" {
		newEntity.State = change.Additions.State
	}

	var oldAttributes RawObject
	if err := convertTypes(&oldAttributes, oldEntity.Attributes); err != nil {
		return Entity[AttributeT]{}, fmt.Errorf("converting the old attributes: %w", err)
	}

	attributes := mergeMaps(
		oldAttributes,
		change.Additions.Attributes,
		change.Removals.Attributes,
	)

	if err := convertTypes(&newEntity.Attributes, attributes); err != nil {
		return Entity[AttributeT]{}, fmt.Errorf("converting the new attributes: %w", err)
	}

	return newEntity, nil
}

func mergeMaps(old, additions RawObject, removals []string) RawObject {
	new := make(RawObject, len(old)+len(additions)-len(removals))
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

func mergeContexts(context, additions Context, removals []string) Context {
	// Adjust context for any additions:
	if additions.ID != nil {
		context.ID = additions.ID
	}
	if additions.UserID != nil {
		context.UserID = additions.UserID
	}
	if additions.ParentID != nil {
		context.ParentID = additions.ParentID
	}

	// Adjust context for any removals:
	for _, key := range removals {
		switch key {
		case "user_id":
			context.UserID = nil
		case "id":
			context.ID = nil
		case "parent_id":
			context.ParentID = nil
		}
	}

	return context
}

// Convert `src` to `dst` (which can be of two different types) by
// serializing to JSON then deserializing. `src` must be something
// that can be passed to `json.Marshal()`, and `dst` must be something
// that can be passed to `json.Unmarshal()` (i.e., typically a
// pointer).
func convertTypes(dst any, src any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("serializing src: %w", err)
	}

	if err := json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("deserializing to dst: %w", err)
	}

	return nil
}
