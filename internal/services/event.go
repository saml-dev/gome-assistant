package services

type Event struct {
	api API
}

/* Public API */

// Fire an event. Takes an event type and an optional map that is sent
// as `event_data`.
func (e Event) Fire(eventType string, eventData ...map[string]any) error {
	if len(eventData) == 0 {
		return e.api.FireEvent(eventType, nil)
	}
	return e.api.FireEvent(eventType, eventData[0])
}
