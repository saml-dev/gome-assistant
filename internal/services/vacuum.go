package services

import (
	"context"

	ga "saml.dev/gome-assistant"
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
func (v Vacuum) CleanSpot(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "clean_spot",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Locate the vacuum cleaner robot.
func (v Vacuum) Locate(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "locate",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Pause the cleaning task.
func (v Vacuum) Pause(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "pause",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Tell the vacuum cleaner to return to its dock.
func (v Vacuum) ReturnToBase(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "return_to_base",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Send a raw command to the vacuum cleaner.
func (v Vacuum) SendCommand(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "send_command",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Set the fan speed of the vacuum cleaner.
func (v Vacuum) SetFanSpeed(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "set_fan_speed",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Start or resume the cleaning task.
func (v Vacuum) Start(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "start",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Start, pause, or resume the cleaning task.
func (v Vacuum) StartPause(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "start_pause",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Stop the current cleaning task.
func (v Vacuum) Stop(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "stop",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Stop the current cleaning task and return to home.
func (v Vacuum) TurnOff(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "turn_off",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Start a new cleaning task.
func (v Vacuum) TurnOn(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := v.service.CallService(
		ctx, "vacuum", "turn_on",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
