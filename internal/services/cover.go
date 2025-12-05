package services

/* Structs */

type Cover struct {
	api API
}

/* Public API */

// Close all or specified cover. Takes an entityId.
func (c Cover) Close(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Close all or specified cover tilt. Takes an entityId.
func (c Cover) CloseTilt(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover_tilt",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Open all or specified cover. Takes an entityId.
func (c Cover) Open(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Open all or specified cover tilt. Takes an entityId.
func (c Cover) OpenTilt(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover_tilt",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Move to specific position all or specified cover. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "set_cover_position",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return c.api.Call(req)
}

// Move to specific position all or specified cover tilt. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Target:  Entity(entityId),
		Domain:  "cover",
		Service: "set_cover_tilt_position",
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return c.api.Call(req)
}

// Stop a cover entity. Takes an entityId.
func (c Cover) Stop(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Stop a cover entity tilt. Takes an entityId.
func (c Cover) StopTilt(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover_tilt",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Toggle a cover open/closed. Takes an entityId.
func (c Cover) Toggle(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}

// Toggle a cover tilt open/closed. Takes an entityId.
func (c Cover) ToggleTilt(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle_cover_tilt",
		Target:  Entity(entityId),
	}
	return c.api.Call(req)
}
