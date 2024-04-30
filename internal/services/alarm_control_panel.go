package services

import (
	"context"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/websocket"
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
func (acp AlarmControlPanel) ArmAway(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_away",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_custom_bypass",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Send the alarm the command for arm home.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_home",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Send the alarm the command for arm night.
func (acp AlarmControlPanel) ArmNight(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_night",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Send the alarm the command for arm vacation.
func (acp AlarmControlPanel) ArmVacation(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_arm_vacation",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Send the alarm the command for disarm.
func (acp AlarmControlPanel) Disarm(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_disarm",
		serviceData, ga.EntityTarget(entityID),
	)
}

// Send the alarm the command for trigger.
func (acp AlarmControlPanel) Trigger(
	entityID string, serviceData any,
) (websocket.Message, error) {
	ctx := context.TODO()
	return acp.service.CallService(
		ctx, "alarm_control_panel", "alarm_trigger",
		serviceData, ga.EntityTarget(entityID),
	)
}
