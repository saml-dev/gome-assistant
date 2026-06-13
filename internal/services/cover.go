package services

import "context"

/* Structs */

type Cover struct {
	api API
}

/* Public API */

// Close all or specified cover. Takes entity IDs.
func (c Cover) Close(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Close all or specified cover tilt. Takes entity IDs.
func (c Cover) CloseTilt(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "close_cover_tilt",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Open all or specified cover. Takes entity IDs.
func (c Cover) Open(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Open all or specified cover tilt. Takes entity IDs.
func (c Cover) OpenTilt(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "open_cover_tilt",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Move to specific position all or specified cover. Takes entity IDs
// and an optional service_data, which must be serializable to a JSON
// object.
func (c Cover) SetPosition(
	ctx context.Context, entityIDs []string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "cover",
		Service:     "set_cover_position",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Move to specific position all or specified cover tilt. Takes an
// entity IDs and an optional service_data, which must be serializable
// to a JSON object.
func (c Cover) SetTiltPosition(
	ctx context.Context, entityIDs []string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Target:      Entities(entityIDs),
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

// Stop a cover entity. Takes entity IDs.
func (c Cover) Stop(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Stop a cover entity tilt. Takes entity IDs.
func (c Cover) StopTilt(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "stop_cover_tilt",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a cover open/closed. Takes entity IDs.
func (c Cover) Toggle(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a cover tilt open/closed. Takes entity IDs.
func (c Cover) ToggleTilt(
	ctx context.Context, entityIDs []string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "cover",
		Service: "toggle_cover_tilt",
		Target:  Entities(entityIDs),
	}

	var result any
	if err := c.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
