package services

import "context"

/* Structs */

type Vacuum struct {
	api API
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes entity IDs.
func (v Vacuum) CleanSpot(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "clean_spot",
		Target:  target,
	}
	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Locate the vacuum cleaner robot.
// Takes entity IDs.
func (v Vacuum) Locate(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "locate",
		Target:  target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Pause the cleaning task.
// Takes entity IDs.
func (v Vacuum) Pause(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "pause",
		Target:  target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Tell the vacuum cleaner to return to its dock.
// Takes entity IDs.
func (v Vacuum) ReturnToBase(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "return_to_base",
		Target:  target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Send a raw command to the vacuum cleaner. Takes entity IDs and an
// optional service_data, which must be serializable to a JSON object.
func (v Vacuum) SendCommand(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "vacuum",
		Service:     "send_command",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Set the fan speed of the vacuum cleaner. Takes entity IDs and an
// optional service_data, which must be serializable to a JSON object.
func (v Vacuum) SetFanSpeed(
	ctx context.Context, target Target, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "vacuum",
		Service:     "set_fan_speed",
		ServiceData: optionalServiceData(serviceData...),
		Target:      target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Start or resume the cleaning task.
// Takes entity IDs.
func (v Vacuum) Start(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "start",
		Target:  target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Start, pause, or resume the cleaning task.
// Takes entity IDs.
func (v Vacuum) StartPause(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "start_pause",
		Target:  target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop the current cleaning task.
// Takes entity IDs.
func (v Vacuum) Stop(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "stop",
		Target:  target,
	}
	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop the current cleaning task and return to home.
// Takes entity IDs.
func (v Vacuum) TurnOff(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "turn_off",
		Target:  target,
	}

	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Start a new cleaning task.
// Takes entity IDs.
func (v Vacuum) TurnOn(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "vacuum",
		Service: "turn_on",
		Target:  target,
	}
	var result any
	if err := v.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
