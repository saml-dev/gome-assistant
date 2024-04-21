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
func (acp AlarmControlPanel) ArmAway(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_away",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Send the alarm the command for arm away.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmWithCustomBypass(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_custom_bypass",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Send the alarm the command for arm home.
// Takes an entityID and an optional
// map that is translated into service_data.
func (acp AlarmControlPanel) ArmHome(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_home",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Send the alarm the command for arm night.
func (acp AlarmControlPanel) ArmNight(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_night",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Send the alarm the command for arm vacation.
func (acp AlarmControlPanel) ArmVacation(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_arm_vacation",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Send the alarm the command for disarm.
func (acp AlarmControlPanel) Disarm(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_disarm",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Send the alarm the command for trigger.
func (acp AlarmControlPanel) Trigger(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "alarm_control_panel",
		Service: "alarm_trigger",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	acp.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}
