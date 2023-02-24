package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type AlarmControlPanel struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// Send the alarm the command for arm away.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmAway(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_away"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}

// Send the alarm the command for arm away.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_custom_bypass"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}

// Send the alarm the command for arm home.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_home"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}

// Send the alarm the command for arm night.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmNight(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_night"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}

// Send the alarm the command for arm vacation.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmVacation(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_vacation"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}

// Send the alarm the command for disarm.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Disarm(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_disarm"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}

// Send the alarm the command for trigger.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Trigger(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_trigger"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	acp.conn.WriteMessage(req, acp.ctx)
}
