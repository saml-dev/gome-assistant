package services

/* Structs */

type Vacuum struct {
	api API
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityId.
func (v Vacuum) CleanSpot(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "clean_spot",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Locate the vacuum cleaner robot.
// Takes an entityId.
func (v Vacuum) Locate(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "locate",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Pause the cleaning task.
// Takes an entityId.
func (v Vacuum) Pause(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "pause",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityId.
func (v Vacuum) ReturnToBase(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "return_to_base",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Send a raw command to the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SendCommand(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "send_command",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return v.api.Call(req)
}

// Set the fan speed of the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SetFanSpeed(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "set_fan_speed",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return v.api.Call(req)
}

// Start or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) Start(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "start",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Start, pause, or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) StartPause(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "start_pause",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Stop the current cleaning task.
// Takes an entityId.
func (v Vacuum) Stop(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "stop",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Stop the current cleaning task and return to home.
// Takes an entityId.
func (v Vacuum) TurnOff(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "turn_off",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}

// Start a new cleaning task.
// Takes an entityId.
func (v Vacuum) TurnOn(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "turn_on",
		Target:  Entity(entityId),
	}
	return v.api.Call(req)
}
