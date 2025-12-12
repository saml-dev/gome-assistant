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
	return c.api.Call(req)
}

// Close all or specified cover tilt. Takes an entityID.
func (c Cover) CloseTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}

// Open all or specified cover. Takes an entityID.
func (c Cover) Open(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}

// Open all or specified cover tilt. Takes an entityID.
func (c Cover) OpenTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}

// Move to specific position all or specified cover. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "set_cover_position",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return c.api.Call(req)
}

// Move to specific position all or specified cover tilt. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Target:  Entity(entityID),
		Domain:  "cover",
		Service: "set_cover_tilt_position",
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return c.api.Call(req)
}

// Stop a cover entity. Takes an entityID.
func (c Cover) Stop(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}

// Stop a cover entity tilt. Takes an entityID.
func (c Cover) StopTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}

// Toggle a cover open/closed. Takes an entityID.
func (c Cover) Toggle(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}

// Toggle a cover tilt open/closed. Takes an entityID.
func (c Cover) ToggleTilt(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle_cover_tilt",
		Target:  Entity(entityID),
	}
	return c.api.Call(req)
}
