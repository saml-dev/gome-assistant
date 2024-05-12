package services

import (
	"context"

	ga "saml.dev/gome-assistant"
)

/* Structs */

type Cover struct {
	service Service
}

func NewCover(service Service) *Cover {
	return &Cover{
		service: service,
	}
}

/* Public API */

// Close all or specified cover. Takes an entityID.
func (c Cover) Close(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "close_cover",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Close all or specified cover tilt. Takes an entityID.
func (c Cover) CloseTilt(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "close_cover_tilt",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Open all or specified cover. Takes an entityID.
func (c Cover) Open(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "open_cover",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Open all or specified cover tilt. Takes an entityID.
func (c Cover) OpenTilt(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "open_cover_tilt",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Move to specific position all or specified cover. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "set_cover_position",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Move to specific position all or specified cover tilt. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(target ga.Target, serviceData any) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "set_cover_tilt_position",
		serviceData, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Stop a cover entity. Takes an entityID.
func (c Cover) Stop(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "stop_cover",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Stop a cover entity tilt. Takes an entityID.
func (c Cover) StopTilt(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "stop_cover_tilt",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Toggle a cover open/closed. Takes an entityID.
func (c Cover) Toggle(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "toggle",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Toggle a cover tilt open/closed. Takes an entityID.
func (c Cover) ToggleTilt(target ga.Target) (any, error) {
	ctx := context.TODO()
	var result any
	err := c.service.CallService(
		ctx, "cover", "toggle_cover_tilt",
		nil, target, &result,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
