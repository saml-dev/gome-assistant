package services

type HomeAssistant struct {
	api API
}

// TurnOn a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return ha.api.WriteMessage(req)
}

// Toggle a Home Assistant entity. Takes an entityId and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "toggle"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return ha.api.WriteMessage(req)
}

func (ha *HomeAssistant) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "homeassistant"
	req.Service = "turn_off"

	return ha.api.WriteMessage(req)
}
