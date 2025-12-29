package services

/* Structs */

type Timer struct {
	api API
}

/* Public API */

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Start(entityID string, duration string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "start",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: Entity(entityID),
	}
	return t.api.CallAndForget(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Change(entityID string, duration string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "change",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: Entity(entityID),
	}
	return t.api.CallAndForget(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerpause
func (t Timer) Pause(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "pause",
		Target:  Entity(entityID),
	}
	return t.api.CallAndForget(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timercancel
func (t Timer) Cancel() error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "cancel",
		Target:  Entity(""),
	}
	return t.api.CallAndForget(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerfinish
func (t Timer) Finish(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "finish",
		Target:  Entity(entityID),
	}
	return t.api.CallAndForget(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerreload
func (t Timer) Reload() error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "reload",
		Target:  Entity(""),
	}
	return t.api.CallAndForget(req)
}
