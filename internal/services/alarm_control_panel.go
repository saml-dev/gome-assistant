package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type AlarmControlPanel struct {
	service Service
}

/* Public API */

func NewAlarmControlPanel(service Service) *AlarmControlPanel {
	return &AlarmControlPanel{
		service: service,
	}
}

// Send the alarm the command for arm away.
func (acp AlarmControlPanel) ArmAway(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_away",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(
	target ga.Target, serviceData any,
) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_custom_bypass",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send the alarm the command for arm home.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_home",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send the alarm the command for arm night.
func (acp AlarmControlPanel) ArmNight(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_night",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send the alarm the command for arm vacation.
func (acp AlarmControlPanel) ArmVacation(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_vacation",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send the alarm the command for disarm.
func (acp AlarmControlPanel) Disarm(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_disarm",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send the alarm the command for trigger.
func (acp AlarmControlPanel) Trigger(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_trigger",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
