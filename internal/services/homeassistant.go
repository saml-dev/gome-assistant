package services

type HomeAssistant struct {
	api API
}

// TurnOn a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_on",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return ha.api.Call(req)
}

// Toggle a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "toggle",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return ha.api.Call(req)
}

func (ha *HomeAssistant) TurnOff(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_off",
		Target:  Entity(entityId),
	}
	return ha.api.Call(req)
}
