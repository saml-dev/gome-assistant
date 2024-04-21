package services

import (
	"context"

	"saml.dev/gome-assistant/internal/websocket"
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
func (c Cover) Close(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "close_cover",
		nil, EntityTarget(entityID),
	)
}

// Close all or specified cover tilt. Takes an entityID.
func (c Cover) CloseTilt(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "close_cover_tilt",
		nil, EntityTarget(entityID),
	)
}

// Open all or specified cover. Takes an entityID.
func (c Cover) Open(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "open_cover",
		nil, EntityTarget(entityID),
	)
}

// Open all or specified cover tilt. Takes an entityID.
func (c Cover) OpenTilt(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "open_cover_tilt",
		nil, EntityTarget(entityID),
	)
}

// Move to specific position all or specified cover. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "set_cover_position",
		serviceData, EntityTarget(entityID),
	)
}

// Move to specific position all or specified cover tilt. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityID string, serviceData any) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "set_cover_tilt_position",
		serviceData, EntityTarget(entityID),
	)
}

// Stop a cover entity. Takes an entityID.
func (c Cover) Stop(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "stop_cover",
		nil, EntityTarget(entityID),
	)
}

// Stop a cover entity tilt. Takes an entityID.
func (c Cover) StopTilt(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "stop_cover_tilt",
		nil, EntityTarget(entityID),
	)
}

// Toggle a cover open/closed. Takes an entityID.
func (c Cover) Toggle(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "toggle",
		nil, EntityTarget(entityID),
	)
}

// Toggle a cover tilt open/closed. Takes an entityID.
func (c Cover) ToggleTilt(entityID string) (websocket.Message, error) {
	ctx := context.TODO()
	return c.service.CallService(
		ctx, "cover", "toggle_cover_tilt",
		nil, EntityTarget(entityID),
	)
}
