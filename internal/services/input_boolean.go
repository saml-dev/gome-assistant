package services

/* Structs */

type InputBoolean struct {
	api API
}

/* Public API */

func (ib InputBoolean) TurnOn(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_on"

	return ib.api.WriteMessage(req)
}

func (ib InputBoolean) Toggle(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "toggle"

	return ib.api.WriteMessage(req)
}

func (ib InputBoolean) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_off"
	return ib.api.WriteMessage(req)
}

func (ib InputBoolean) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "input_boolean"
	req.Service = "reload"
	return ib.api.WriteMessage(req)
}
