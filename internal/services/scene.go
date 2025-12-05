package services

/* Structs */

type Scene struct {
	api API
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "scene",
		Service: "apply",
		Target:  Entity(""),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return s.api.Call(req)
}

// Create a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) Create(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "scene",
		Service: "create",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return s.api.Call(req)
}

// Reload the scenes.
func (s Scene) Reload() error {
	req := BaseServiceRequest{
		Domain:  "scene",
		Service: "reload",
		Target:  Entity(""),
	}
	return s.api.Call(req)
}

// TurnOn a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "scene",
		Service: "turn_on",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return s.api.Call(req)
}
