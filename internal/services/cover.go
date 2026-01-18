package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type Cover struct {
	api API
}

/* Public API */

// Close all or specified cover. Takes an entityID.
func (c Cover) Close(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Close all or specified cover tilt. Takes an entityID.
func (c Cover) CloseTilt(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover_tilt",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Open all or specified cover. Takes an entityID.
func (c Cover) Open(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Open all or specified cover tilt. Takes an entityID.
func (c Cover) OpenTilt(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover_tilt",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Move to specific position all or specified cover. Takes an entityID
// and an optional service_data, which must be serializable to a JSON
// object.
func (c Cover) SetPosition(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "cover",
		Service:     "set_cover_position",
		ServiceData: optionalServiceData(serviceData...),
		Target:      message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Move to specific position all or specified cover tilt. Takes an
// entityID and an optional service_data, which must be serializable
// to a JSON object.
func (c Cover) SetTiltPosition(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Target:      message.Entity(entityID),
		Domain:      "cover",
		ServiceData: optionalServiceData(serviceData...),
		Service:     "set_cover_tilt_position",
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop a cover entity. Takes an entityID.
func (c Cover) Stop(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop a cover entity tilt. Takes an entityID.
func (c Cover) StopTilt(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover_tilt",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a cover open/closed. Takes an entityID.
func (c Cover) Toggle(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a cover tilt open/closed. Takes an entityID.
func (c Cover) ToggleTilt(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle_cover_tilt",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
