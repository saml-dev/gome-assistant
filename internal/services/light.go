package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Light struct {
	conn *websocket.Conn
}

func NewLight(conn *websocket.Conn) *Light {
	return &Light{
		conn: conn,
	}
}

/* Public API */

// TurnOn a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) TurnOn(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(l.conn, entityId)
	req.Domain = "light"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Toggle a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) Toggle(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(l.conn, entityId)
	req.Domain = "light"
	req.Service = "toggle"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (l Light) TurnOff(entityId string) {
	req := NewBaseServiceRequest(l.conn, entityId)
	req.Domain = "light"
	req.Service = "turn_off"

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
