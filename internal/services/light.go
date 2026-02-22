package services

/* Structs */

type Light struct {
	api API
}

/* Public API */

// TurnOn a light entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (l Light) TurnOn(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return l.api.CallAndForget(req)
}

// Toggle a light entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (l Light) Toggle(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "light",
		Service:     "toggle",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return l.api.CallAndForget(req)
}

func (l Light) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return l.api.CallAndForget(req)
}
