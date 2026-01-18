package message

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
	TimeFired TimeStamp    `json:"time_fired"`
	Context   EventContext `json:"context"`
}

// EventOrigin represents the origin of an event.
type EventOrigin string

type EventContext struct {
	ID       string `json:"id,omitzero"`
	ParentID string `json:"parent_id,omitzero"`
	UserID   string `json:"user_id,omitzero"`
}

type ZWaveJSValueNotificationEventMessage EventMessage[ZWaveJSValueNotificationData]

type ZWaveJSValueNotificationData struct {
	Domain           string `json:"domain"`
	NodeID           int    `json:"node_id"`
	HomeID           int64  `json:"home_id"`
	Endpoint         int    `json:"endpoint"`
	DeviceID         string `json:"device_id"`
	CommandClass     int    `json:"command_class"`
	CommandClassName string `json:"command_class_name"`
	Label            string `json:"label"`
	Property         string `json:"property"`
	PropertyName     string `json:"property_name"`
	PropertyKey      string `json:"property_key"`
	PropertyKeyName  string `json:"property_key_name"`
	Value            string `json:"value"`
	ValueRaw         int    `json:"value_raw"`
}
