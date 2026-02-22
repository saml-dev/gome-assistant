package services

type HomeAssistant struct {
	api API
}

// TurnOn a Home Assistant entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (ha *HomeAssistant) TurnOn(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "homeassistant",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return ha.api.CallAndForget(req)
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (ha *HomeAssistant) Toggle(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "homeassistant",
		Service:     "toggle",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return ha.api.CallAndForget(req)
}

func (ha *HomeAssistant) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return ha.api.CallAndForget(req)
}
