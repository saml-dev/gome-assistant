package message

import "time"

// EventMessage represents a full event message. The `DataT` type
// parameter specifies the type of the `Event.Data` field. See
// the Home Assistant [docs] for more information.
//
// [docs]: https://developers.home-assistant.io/docs/api/websocket#subscribe-to-events
type EventMessage[DataT any] struct {
	BaseMessage
	Event Event[DataT] `json:"event"`
}

// Event represents the `event` field of an event message. The `DataT`
// type parameter specifies the type of the `Data` field.
type Event[DataT any] struct {
	EventType string       `json:"event_type"`
	Data      DataT        `json:"data"`
	Origin    EventOrigin  `json:"origin"`
	TimeFired time.Time    `json:"time_fired"`
	Context   EventContext `json:"context"`
}

// EventOrigin represents the origin of an event.
type EventOrigin string

type EventContext struct {
	ID       string `json:"id,omitzero"`
	ParentID string `json:"parent_id,omitzero"`
	UserID   string `json:"user_id,omitzero"`
}
