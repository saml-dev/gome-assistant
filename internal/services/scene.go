package services

/* Structs */

type Scene struct {
	api API
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData ...map[string]any) error {
	req := NewBaseServiceRequest("")
	req.Domain = "scene"
	req.Service = "apply"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return s.api.WriteMessage(req)
}

// Create a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) Create(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "scene"
	req.Service = "create"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return s.api.WriteMessage(req)
}

// Reload the scenes.
func (s Scene) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "scene"
	req.Service = "reload"

	return s.api.WriteMessage(req)
}

// TurnOn a scene entity. Takes an entityId and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "scene"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return s.api.WriteMessage(req)
}
