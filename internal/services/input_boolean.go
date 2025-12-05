package services

/* Structs */

type InputBoolean struct {
	api API
}

/* Public API */

func (ib InputBoolean) TurnOn(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_on",
		Target:  Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputBoolean) Toggle(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "toggle",
		Target:  Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputBoolean) TurnOff(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_off",
		Target:  Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputBoolean) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "reload",
	}
	return ib.api.Call(req)
}
