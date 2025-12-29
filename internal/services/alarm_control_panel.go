package services

/* Structs */

type AlarmControlPanel struct {
	api API
}

/* Public API */

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmAway(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_away",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_custom_bypass",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm home.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_home",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm night.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmNight(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_night",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm vacation.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmVacation(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_vacation",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for disarm.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Disarm(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_disarm",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for trigger.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Trigger(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_trigger",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.CallAndForget(req)
}
