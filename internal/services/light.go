package services

/* Structs */

type Light struct {
	api API
}

/* Public API */

// TurnOn a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_on",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return l.api.Call(req)
}

// Toggle a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) Toggle(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "toggle",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return l.api.Call(req)
}

func (l Light) TurnOff(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target:  Entity(entityId),
	}
	return l.api.Call(req)
}
