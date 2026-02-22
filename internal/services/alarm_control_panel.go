package services

/* Structs */

type AlarmControlPanel struct {
	api API
}

/* Public API */

// Send the alarm the command for arm away. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (acp AlarmControlPanel) ArmAway(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_away",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm away. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_custom_bypass",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm home. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (acp AlarmControlPanel) ArmHome(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_home",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm night. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (acp AlarmControlPanel) ArmNight(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_night",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for arm vacation. Takes an entityID and
// an optional service_data, which must be serializable to a JSON
// object.
func (acp AlarmControlPanel) ArmVacation(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_vacation",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for disarm. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (acp AlarmControlPanel) Disarm(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_disarm",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}

// Send the alarm the command for trigger. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (acp AlarmControlPanel) Trigger(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "alarm_control_panel",
		Service:     "alarm_trigger",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return acp.api.CallAndForget(req)
}
