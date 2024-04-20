package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Vacuum struct {
	conn *websocket.Conn
}

func NewVacuum(conn *websocket.Conn) *Vacuum {
	return &Vacuum{
		conn: conn,
	}
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
func (v Vacuum) CleanSpot(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "clean_spot",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Locate the vacuum cleaner robot.
func (v Vacuum) Locate(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "locate",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Pause the cleaning task.
func (v Vacuum) Pause(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "pause",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Tell the vacuum cleaner to return to its dock.
func (v Vacuum) ReturnToBase(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "return_to_base",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Send a raw command to the vacuum cleaner.
func (v Vacuum) SendCommand(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "send_command",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Set the fan speed of the vacuum cleaner.
func (v Vacuum) SetFanSpeed(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "set_fan_speed",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Start or resume the cleaning task.
func (v Vacuum) Start(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "start",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Start, pause, or resume the cleaning task.
func (v Vacuum) StartPause(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "start_pause",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Stop the current cleaning task.
func (v Vacuum) Stop(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "stop",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Stop the current cleaning task and return to home.
func (v Vacuum) TurnOff(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "turn_off",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Start a new cleaning task.
func (v Vacuum) TurnOn(entityID string) {
	req := CallServiceRequest{
		Domain:  "vacuum",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
	}

	v.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
