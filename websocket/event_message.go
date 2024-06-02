package websocket

import "time"

type BaseEvent struct {
	EventType string    `json:"event_type"`
	Origin    string    `json:"origin"`
	TimeFired time.Time `json:"time_fired"`
	Context   Context   `json:"context"`
}

type Event struct {
	BaseEvent
	RawData RawMessage `json:"data"`
}

type EventMessage struct {
	BaseMessage
	Event Event `json:"event"`
}
