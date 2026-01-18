package services

type Event struct {
	api API
}

/* Public API */

// Fire an event. Takes an event type and an optional map that is sent
// as `event_data`.
func (e Event) Fire(eventType string, eventData any) error {
	return e.api.FireEvent(eventType, eventData)
}
