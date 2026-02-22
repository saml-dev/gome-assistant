package services

/* Structs */

type Cover struct {
	api API
}

/* Public API */

// Close all or specified cover. Takes an entityID.
func (c Cover) Close(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Close all or specified cover tilt. Takes an entityID.
func (c Cover) CloseTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Open all or specified cover. Takes an entityID.
func (c Cover) Open(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Open all or specified cover tilt. Takes an entityID.
func (c Cover) OpenTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Move to specific position all or specified cover. Takes an entityID
// and an optional service_data, which must be serializable to a JSON
// object.
func (c Cover) SetPosition(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "cover",
		Service:     "set_cover_position",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return c.api.CallAndForget(req)
}

// Move to specific position all or specified cover tilt. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (c Cover) SetTiltPosition(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Target:      Entity(entityID),
		Domain:      "cover",
		ServiceData: optionalServiceData(serviceData...),
		Service:     "set_cover_tilt_position",
	}

	return c.api.CallAndForget(req)
}

// Stop a cover entity. Takes an entityID.
func (c Cover) Stop(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Stop a cover entity tilt. Takes an entityID.
func (c Cover) StopTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Toggle a cover open/closed. Takes an entityID.
func (c Cover) Toggle(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}

// Toggle a cover tilt open/closed. Takes an entityID.
func (c Cover) ToggleTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.CallAndForget(req)
}
