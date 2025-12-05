package services

/* Structs */

type AlarmControlPanel struct {
	api API
}

/* Public API */

// Send the alarm the command for arm away.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmAway(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_away",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}

// Send the alarm the command for arm away.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_custom_bypass",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}

// Send the alarm the command for arm home.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_home",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}

// Send the alarm the command for arm night.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmNight(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_night",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}

// Send the alarm the command for arm vacation.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmVacation(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_vacation",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}

// Send the alarm the command for disarm.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Disarm(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_disarm",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}

// Send the alarm the command for trigger.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Trigger(entityId string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_trigger",
		Target:  Entity(entityId),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return acp.api.Call(req)
}
