package services

/* Structs */

type Switch struct {
	api API
}

/* Public API */

func (s Switch) TurnOn(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "switch",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	return s.api.CallAndForget(req)
}

func (s Switch) Toggle(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "switch",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	return s.api.CallAndForget(req)
}

func (s Switch) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "switch",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return s.api.CallAndForget(req)
}
