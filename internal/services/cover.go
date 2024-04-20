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

// Close all or specified cover. Takes an entityID.
func (c Cover) Close(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "close_cover",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Close all or specified cover tilt. Takes an entityID.
func (c Cover) CloseTilt(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "close_cover_tilt",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Open all or specified cover. Takes an entityID.
func (c Cover) Open(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "open_cover",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Open all or specified cover tilt. Takes an entityID.
func (c Cover) OpenTilt(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "open_cover_tilt",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Move to specific position all or specified cover. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetPosition(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "set_cover_position",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Move to specific position all or specified cover tilt. Takes an entityID and an optional
// map that is translated into service_data.
func (c Cover) SetTiltPosition(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "set_cover_tilt_position",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Stop a cover entity. Takes an entityID.
func (c Cover) Stop(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "stop_cover",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Stop a cover entity tilt. Takes an entityID.
func (c Cover) StopTilt(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "stop_cover_tilt",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a cover open/closed. Takes an entityID.
func (c Cover) Toggle(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a cover tilt open/closed. Takes an entityID.
func (c Cover) ToggleTilt(entityID string) {
	req := CallServiceRequest{
		Domain:  "cover",
		Service: "toggle_cover_tilt",
		Target: Target{
			EntityID: entityID,
		},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
