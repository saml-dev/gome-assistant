package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Light struct {
	conn *ws.WebsocketConn
}

/* Public API */

// TurnOn a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) TurnOn(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_on"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return l.conn.WriteMessage(req)
}

// Toggle a light entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Light) Toggle(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "toggle"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return l.conn.WriteMessage(req)
}

func (l Light) TurnOff(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "light"
	req.Service = "turn_off"
	return l.conn.WriteMessage(req)
}
