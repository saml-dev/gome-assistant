package services

/* Structs */

type Script struct {
	api API
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "reload",
		Target:  Entity(entityID),
	}
	return s.api.Call(req)
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	return s.api.Call(req)
}

// TurnOff a script that was created in the HA UI.
func (s Script) TurnOff() error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "turn_off",
		Target:  Entity(""),
	}
	return s.api.Call(req)
}

// TurnOn a script that was created in the HA UI.
func (s Script) TurnOn(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	return s.api.Call(req)
}
