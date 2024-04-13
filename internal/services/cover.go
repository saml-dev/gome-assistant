package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Cover struct {
	conn *ws.WebsocketConn
}

/* Public API */

// Close all or specified cover. Takes an entityId.
func (c Cover) Close(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "close_cover"

	return c.conn.WriteMessage(req)
}

// Close all or specified cover tilt. Takes an entityId.
func (c Cover) CloseTilt(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "close_cover_tilt"

	return c.conn.WriteMessage(req)
}

// Open all or specified cover. Takes an entityId.
func (c Cover) Open(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "open_cover"

	return c.conn.WriteMessage(req)
}

// Open all or specified cover tilt. Takes an entityId.
func (c Cover) OpenTilt(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "open_cover_tilt"

	return c.conn.WriteMessage(req)
}

// Move to specific position all or specified cover. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "set_cover_position"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return c.conn.WriteMessage(req)
}

// Move to specific position all or specified cover tilt. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "set_cover_tilt_position"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return c.conn.WriteMessage(req)
}

// Stop a cover entity. Takes an entityId.
func (c Cover) Stop(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "stop_cover"

	return c.conn.WriteMessage(req)
}

// Stop a cover entity tilt. Takes an entityId.
func (c Cover) StopTilt(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "stop_cover_tilt"

	return c.conn.WriteMessage(req)
}

// Toggle a cover open/closed. Takes an entityId.
func (c Cover) Toggle(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "toggle"

	return c.conn.WriteMessage(req)
}

// Toggle a cover tilt open/closed. Takes an entityId.
func (c Cover) ToggleTilt(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "cover"
	req.Service = "toggle_cover_tilt"

	return c.conn.WriteMessage(req)
}
