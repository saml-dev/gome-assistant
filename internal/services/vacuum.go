package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Vacuum struct {
	conn *websocket.Conn
}

func NewVacuum(conn *websocket.Conn) *Vacuum {
	return &Vacuum{
		conn: conn,
	}
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityId.
func (v Vacuum) CleanSpot(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "clean_spot"

	v.conn.WriteMessage(req)
}

// Locate the vacuum cleaner robot.
// Takes an entityId.
func (v Vacuum) Locate(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "locate"

	v.conn.WriteMessage(req)
}

// Pause the cleaning task.
// Takes an entityId.
func (v Vacuum) Pause(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "pause"

	v.conn.WriteMessage(req)
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityId.
func (v Vacuum) ReturnToBase(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "return_to_base"

	v.conn.WriteMessage(req)
}

// Send a raw command to the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SendCommand(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "send_command"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	v.conn.WriteMessage(req)
}

// Set the fan speed of the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SetFanSpeed(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "set_fan_speed"

	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	v.conn.WriteMessage(req)
}

// Start or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) Start(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "start"

	v.conn.WriteMessage(req)
}

// Start, pause, or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) StartPause(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "start_pause"

	v.conn.WriteMessage(req)
}

// Stop the current cleaning task.
// Takes an entityId.
func (v Vacuum) Stop(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "stop"

	v.conn.WriteMessage(req)
}

// Stop the current cleaning task and return to home.
// Takes an entityId.
func (v Vacuum) TurnOff(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "turn_off"

	v.conn.WriteMessage(req)
}

// Start a new cleaning task.
// Takes an entityId.
func (v Vacuum) TurnOn(entityId string) {
	req := NewBaseServiceRequest(v.conn, entityId)
	req.Domain = "vacuum"
	req.Service = "turn_on"

	v.conn.WriteMessage(req)
}
