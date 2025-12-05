package services

/* Structs */

type Light struct {
	api API
}

/* Public API */

// TurnOn a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return l.api.WriteMessage(req)
}

// Toggle a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) Toggle(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "toggle"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return l.api.WriteMessage(req)
}

func (l Light) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_off"
	return l.api.WriteMessage(req)
}
