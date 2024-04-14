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
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmAway(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_away"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm away.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_custom_bypass"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm home.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_home"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm night.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmNight(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_night"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for arm vacation.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmVacation(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_arm_vacation"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for disarm.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Disarm(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_disarm"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send the alarm the command for trigger.
// Takes an entityId and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) Trigger(entityId string, serviceData map[string]any) {
	req := NewBaseServiceRequest(acp.conn, entityId)
	req.Domain = "alarm_control_panel"
	req.Service = "alarm_trigger"
	req.ServiceData = serviceData

	acp.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
