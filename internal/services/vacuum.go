package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Vacuum struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityId.
func (v Vacuum) CleanSpot(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "clean_spot"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Locate the vacuum cleaner robot.
// Takes an entityId.
func (v Vacuum) Locate(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "locate"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Pause the cleaning task.
// Takes an entityId.
func (v Vacuum) Pause(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "pause"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityId.
func (v Vacuum) ReturnToBase(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "return_to_base"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Send a raw command to the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SendCommand(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "send_command"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Set the fan speed of the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SetFanSpeed(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "set_fan_speed"

	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Start or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) Start(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "start"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Start, pause, or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) StartPause(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "start_pause"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Stop the current cleaning task.
// Takes an entityId.
func (v Vacuum) Stop(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "stop"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Stop the current cleaning task and return to home.
// Takes an entityId.
func (v Vacuum) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "turn_off"

	ws.WriteMessage(req, v.conn, v.ctx)
}

// Start a new cleaning task.
// Takes an entityId.
func (v Vacuum) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "turn_on"

	ws.WriteMessage(req, v.conn, v.ctx)
}
