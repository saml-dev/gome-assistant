package services

/* Structs */

type Vacuum struct {
	api API
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityID.
func (v Vacuum) CleanSpot(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "clean_spot",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Locate the vacuum cleaner robot.
// Takes an entityID.
func (v Vacuum) Locate(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "locate",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Pause the cleaning task.
// Takes an entityID.
func (v Vacuum) Pause(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "pause",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityID.
func (v Vacuum) ReturnToBase(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "return_to_base",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Send a raw command to the vacuum cleaner. Takes an entityID and an optional
// map that is translated into service_data.
func (v Vacuum) SendCommand(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "send_command",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return v.api.CallAndForget(req)
}

// Set the fan speed of the vacuum cleaner. Takes an entityID and an optional
// map that is translated into service_data.
func (v Vacuum) SetFanSpeed(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "set_fan_speed",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return v.api.CallAndForget(req)
}

// Start or resume the cleaning task.
// Takes an entityID.
func (v Vacuum) Start(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "start",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Start, pause, or resume the cleaning task.
// Takes an entityID.
func (v Vacuum) StartPause(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "start_pause",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Stop the current cleaning task.
// Takes an entityID.
func (v Vacuum) Stop(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "stop",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Stop the current cleaning task and return to home.
// Takes an entityID.
func (v Vacuum) TurnOff(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "turn_off",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}

// Start a new cleaning task.
// Takes an entityID.
func (v Vacuum) TurnOn(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "turn_on",
		Target:  Entity(entityID),
	}
	return v.api.CallAndForget(req)
}
