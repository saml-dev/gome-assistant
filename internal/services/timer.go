package services

/* Structs */

type Timer struct {
	api API
}

/* Public API */

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Start(entityId string, duration string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "start",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: Entity(entityId),
	}
	return t.api.Call(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Change(entityId string, duration string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "change",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: Entity(entityId),
	}
	return t.api.Call(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerpause
func (t Timer) Pause(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "pause",
		Target:  Entity(entityId),
	}
	return t.api.Call(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timercancel
func (t Timer) Cancel() error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "cancel",
		Target:  Entity(""),
	}
	return t.api.Call(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerfinish
func (t Timer) Finish(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "finish",
		Target:  Entity(entityId),
	}
	return t.api.Call(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerreload
func (t Timer) Reload() error {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "reload",
		Target:  Entity(""),
	}
	return t.api.Call(req)
}
