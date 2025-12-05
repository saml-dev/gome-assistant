package services

/* Structs */

type Script struct {
	api API
}

/* Public API */

// Reload a script that was created in the HA UI.
func (s Script) Reload(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "reload"

	return s.api.WriteMessage(req)
}

// Toggle a script that was created in the HA UI.
func (s Script) Toggle(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "toggle"

	return s.api.WriteMessage(req)
}

// TurnOff a script that was created in the HA UI.
func (s Script) TurnOff() error {
	req := NewBaseServiceRequest("")
	req.Domain = "script"
	req.Service = "turn_off"

	return s.api.WriteMessage(req)
}

// TurnOn a script that was created in the HA UI.
func (s Script) TurnOn(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "script"
	req.Service = "turn_on"

	return s.api.WriteMessage(req)
}
