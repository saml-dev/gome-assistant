package services

/* Structs */

type Switch struct {
	api API
}

/* Public API */

func (s Switch) TurnOn(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_on"

	return s.api.WriteMessage(req)
}

func (s Switch) Toggle(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	return s.api.WriteMessage(req)
}

func (s Switch) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_off"

	return s.api.WriteMessage(req)
}
