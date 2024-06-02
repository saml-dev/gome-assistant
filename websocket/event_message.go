package websocket

import "time"

type EventContext struct {
	ID       *string `json:"id"`
	UserID   *string `json:"user_id"`
	ParentID *string `json:"parent_id"`
}

type BaseEvent struct {
	EventType string       `json:"event_type"`
	Origin    string       `json:"origin"`
	TimeFired time.Time    `json:"time_fired"`
	Context   EventContext `json:"context"`
}

type Event struct {
	BaseEvent
	RawData RawMessage `json:"data"`
}

type EventMessage struct {
	BaseMessage
	Event Event `json:"event"`
}
