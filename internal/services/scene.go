package services

/* Structs */

type Scene struct {
	api API
}

/* Public API */

// Apply a scene. Takes an optional service_data, which must be
// serializable to a JSON object.
func (s Scene) Apply(serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "apply",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(""),
	}

	return s.api.CallAndForget(req)
}

// Create a scene entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (s Scene) Create(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "create",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
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
// service_data, which must be serializable to a JSON object.
func (s Scene) TurnOn(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "scene",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return s.api.CallAndForget(req)
}
