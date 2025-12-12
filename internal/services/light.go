package services

/* Structs */

type Light struct {
	api API
}

/* Public API */

// TurnOn a light entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Light) TurnOn(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return l.api.Call(req)
}

// Toggle a light entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Light) Toggle(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return l.api.Call(req)
}

func (l Light) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "light",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return l.api.Call(req)
}
