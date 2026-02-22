package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type Vacuum struct {
	api API
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityID.
func (v Vacuum) CleanSpot(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "clean_spot",
		Target:  message.Entity(entityID),
	}
	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Locate the vacuum cleaner robot.
// Takes an entityID.
func (v Vacuum) Locate(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "locate",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Pause the cleaning task.
// Takes an entityID.
func (v Vacuum) Pause(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "pause",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityID.
func (v Vacuum) ReturnToBase(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "return_to_base",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send a raw command to the vacuum cleaner. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (v Vacuum) SendCommand(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "vacuum",
		Service:     "send_command",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set the fan speed of the vacuum cleaner. Takes an entityID and an
// optional service_data, which must be serializable to a JSON object.
func (v Vacuum) SetFanSpeed(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "vacuum",
		Service:     "set_fan_speed",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Start or resume the cleaning task.
// Takes an entityID.
func (v Vacuum) Start(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "start",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Start, pause, or resume the cleaning task.
// Takes an entityID.
func (v Vacuum) StartPause(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "start_pause",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop the current cleaning task.
// Takes an entityID.
func (v Vacuum) Stop(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "stop",
		Target:  message.Entity(entityID),
	}
	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop the current cleaning task and return to home.
// Takes an entityID.
func (v Vacuum) TurnOff(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "turn_off",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Start a new cleaning task.
// Takes an entityID.
func (v Vacuum) TurnOn(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "vacuum",
		Service: "turn_on",
		Target:  message.Entity(entityID),
	}
	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
