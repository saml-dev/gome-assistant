package services

/* Structs */

type Script struct {
	api API
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "reload",
		Target:  Entity(entityId),
	}
	return s.api.Call(req)
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "toggle",
		Target:  Entity(entityId),
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
func (s Script) TurnOn(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "script",
		Service: "turn_on",
		Target:  Entity(entityId),
	}
	return s.api.Call(req)
}
