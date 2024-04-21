package services

import (
	"context"

	"saml.dev/gome-assistant/websocket"
)

/* Structs */

type Vacuum struct {
	service Service
}

func NewVacuum(service Service) *Vacuum {
	return &Vacuum{
		service: service,
	}
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
func (v Vacuum) CleanSpot(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "clean_spot",
		nil, EntityTarget(entityID),
	)
}

// Locate the vacuum cleaner robot.
func (v Vacuum) Locate(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "locate",
		nil, EntityTarget(entityID),
	)
}

// Pause the cleaning task.
func (v Vacuum) Pause(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "pause",
		nil, EntityTarget(entityID),
	)
}

// Tell the vacuum cleaner to return to its dock.
func (v Vacuum) ReturnToBase(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "return_to_base",
		nil, EntityTarget(entityID),
	)
}

// Send a raw command to the vacuum cleaner.
func (v Vacuum) SendCommand(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "send_command",
		serviceData, EntityTarget(entityID),
	)
}

// Set the fan speed of the vacuum cleaner.
func (v Vacuum) SetFanSpeed(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "set_fan_speed",
		serviceData, EntityTarget(entityID),
	)
}

// Start or resume the cleaning task.
func (v Vacuum) Start(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "start",
		nil, EntityTarget(entityID),
	)
}

// Start, pause, or resume the cleaning task.
func (v Vacuum) StartPause(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "start_pause",
		nil, EntityTarget(entityID),
	)
}

// Stop the current cleaning task.
func (v Vacuum) Stop(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "stop",
		nil, EntityTarget(entityID),
	)
}

// Stop the current cleaning task and return to home.
func (v Vacuum) TurnOff(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "turn_off",
		nil, EntityTarget(entityID),
	)
}

// Start a new cleaning task.
func (v Vacuum) TurnOn(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return v.service.CallService(
		ctx, "vacuum", "turn_on",
		nil, EntityTarget(entityID),
	)
}
