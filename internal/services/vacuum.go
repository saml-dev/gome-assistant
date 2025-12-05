package services

/* Structs */

type Vacuum struct {
	api API
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityId.
func (v Vacuum) CleanSpot(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "clean_spot"

	return v.api.WriteMessage(req)
}

// Locate the vacuum cleaner robot.
// Takes an entityId.
func (v Vacuum) Locate(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "locate"

	return v.api.WriteMessage(req)
}

// Pause the cleaning task.
// Takes an entityId.
func (v Vacuum) Pause(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "pause"

	return v.api.WriteMessage(req)
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityId.
func (v Vacuum) ReturnToBase(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "return_to_base"

	return v.api.WriteMessage(req)
}

// Send a raw command to the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SendCommand(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "send_command"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return v.api.WriteMessage(req)
}

// Set the fan speed of the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SetFanSpeed(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "set_fan_speed"

	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return v.api.WriteMessage(req)
}

// Start or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) Start(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "start"

	return v.api.WriteMessage(req)
}

// Start, pause, or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) StartPause(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "start_pause"

	return v.api.WriteMessage(req)
}

// Stop the current cleaning task.
// Takes an entityId.
func (v Vacuum) Stop(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "stop"

	return v.api.WriteMessage(req)
}

// Stop the current cleaning task and return to home.
// Takes an entityId.
func (v Vacuum) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "turn_off"

	return v.api.WriteMessage(req)
}

// Start a new cleaning task.
// Takes an entityId.
func (v Vacuum) TurnOn(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "turn_on"

	return v.api.WriteMessage(req)
}
