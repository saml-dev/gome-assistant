package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Cover struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// Close all or specified cover. Takes an entityId.
func (c Cover) Close(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "close_cover"

	c.conn.WriteMessage(req, c.ctx)
}

// Close all or specified cover tilt. Takes an entityId.
func (c Cover) CloseTilt(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "close_cover_tilt"

	c.conn.WriteMessage(req, c.ctx)
}

// Open all or specified cover. Takes an entityId.
func (c Cover) Open(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "open_cover"

	c.conn.WriteMessage(req, c.ctx)
}

// Open all or specified cover tilt. Takes an entityId.
func (c Cover) OpenTilt(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "open_cover_tilt"

	c.conn.WriteMessage(req, c.ctx)
}

// Move to specific position all or specified cover. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "set_cover_position"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	c.conn.WriteMessage(req, c.ctx)
}

// Move to specific position all or specified cover tilt. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "set_cover_tilt_position"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	c.conn.WriteMessage(req, c.ctx)
}

// Stop a cover entity. Takes an entityId.
func (c Cover) Stop(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "stop_cover"

	c.conn.WriteMessage(req, c.ctx)
}

// Stop a cover entity tilt. Takes an entityId.
func (c Cover) StopTilt(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "stop_cover_tilt"

	c.conn.WriteMessage(req, c.ctx)
}

// Toggle a cover open/closed. Takes an entityId.
func (c Cover) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "toggle"

	c.conn.WriteMessage(req, c.ctx)
}

// Toggle a cover tilt open/closed. Takes an entityId.
func (c Cover) ToggleTilt(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "toggle_cover_tilt"

	c.conn.WriteMessage(req, c.ctx)
}
