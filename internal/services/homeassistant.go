package services

type HomeAssistant struct {
	api API
}

// TurnOn a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return ha.api.CallAndForget(req)
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
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
