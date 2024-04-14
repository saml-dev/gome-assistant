package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Cover struct {
	conn *websocket.Conn
}

func NewCover(conn *websocket.Conn) *Cover {
	return &Cover{
		conn: conn,
	}
}

/* Public API */

// Close all or specified cover. Takes an entityId.
func (c Cover) Close(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "close_cover"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Close all or specified cover tilt. Takes an entityId.
func (c Cover) CloseTilt(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "close_cover_tilt"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Open all or specified cover. Takes an entityId.
func (c Cover) Open(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "open_cover"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Open all or specified cover tilt. Takes an entityId.
func (c Cover) OpenTilt(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "open_cover_tilt"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Move to specific position all or specified cover. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "set_cover_position"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Move to specific position all or specified cover tilt. Takes an entityId and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "set_cover_tilt_position"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Stop a cover entity. Takes an entityId.
func (c Cover) Stop(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "stop_cover"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Stop a cover entity tilt. Takes an entityId.
func (c Cover) StopTilt(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "stop_cover_tilt"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a cover open/closed. Takes an entityId.
func (c Cover) Toggle(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "toggle"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a cover tilt open/closed. Takes an entityId.
func (c Cover) ToggleTilt(entityId string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "cover"
	req.Service = "toggle_cover_tilt"

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
