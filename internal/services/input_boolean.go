package services

/* Structs */

type InputBoolean struct {
	api API
}

/* Public API */

func (ib InputBoolean) TurnOn(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	return ib.api.CallAndForget(req)
}

func (ib InputBoolean) Toggle(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	return ib.api.CallAndForget(req)
}

func (ib InputBoolean) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return ib.api.CallAndForget(req)
}

func (ib InputBoolean) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_boolean",
		Service: "reload",
	}
	return ib.api.CallAndForget(req)
}
