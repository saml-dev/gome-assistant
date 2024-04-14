package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type AlarmControlPanel struct {
	conn *websocket.Conn
}

/* Public API */

func NewAlarmControlPanel(conn *websocket.Conn) *AlarmControlPanel {
	return &AlarmControlPanel{
		conn: conn,
	}
}

// Send the alarm the command for arm away.
func (acp AlarmControlPanel) ArmAway(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_away"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_custom_bypass"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm home.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_home"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm night.
func (acp AlarmControlPanel) ArmNight(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_night"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm vacation.
func (acp AlarmControlPanel) ArmVacation(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_vacation"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for disarm.
func (acp AlarmControlPanel) Disarm(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_disarm"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for trigger.
func (acp AlarmControlPanel) Trigger(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{}
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_trigger"
	req.Target.EntityID = entityID
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
