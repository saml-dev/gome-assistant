package services

/* Structs */

type Scene struct {
	api API
}

/* Public API */

// Apply a scene. Takes map that is translated into service_data.
func (s Scene) Apply(serviceData any) error {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "apply",
		Target:      Entity(""),
		ServiceData: serviceData,
	}

	return s.api.CallAndForget(req)
}

// Create a scene entity. Takes an entityID and an optional
// map that is translated into service_data.
func (s Scene) Create(entityID string, serviceData any) error {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "create",
		Target:      Entity(entityID),
		ServiceData: serviceData,
	}

	return s.api.CallAndForget(req)
}

// Reload the scenes.
func (s Scene) Reload() error {
	req := BaseServiceRequest{
		Domain:  "scene",
		Service: "reload",
		Target:  Entity(""),
	}

	return s.api.CallAndForget(req)
}

// TurnOn a scene entity. Takes an entityID and an optional
// map that is translated into service_data.
func (s Scene) TurnOn(entityID string, serviceData any) error {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "turn_on",
		Target:      Entity(entityID),
		ServiceData: serviceData,
	}

	return s.api.CallAndForget(req)
}
