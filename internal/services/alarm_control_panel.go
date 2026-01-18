package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type AlarmControlPanel struct {
	api API
}

/* Public API */

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmAway(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_away",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_custom_bypass",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the alarm the command for arm home.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_home",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the alarm the command for arm night.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmNight(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_night",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the alarm the command for arm vacation.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmVacation(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_arm_vacation",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the alarm the command for disarm.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Disarm(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_disarm",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send the alarm the command for trigger.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Trigger(
	ctx context.Context, entityID string, serviceData any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "alarm_control_panel",
		Service:     "alarm_trigger",
		Target:      message.Entity(entityID),
		ServiceData: serviceData,
	}

	var result any
	if err := acp.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
