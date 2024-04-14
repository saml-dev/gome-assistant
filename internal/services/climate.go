package services

import (
	"saml.dev/gome-assistant/internal/websocket"
	"saml.dev/gome-assistant/types"
)

/* Structs */

type Climate struct {
	conn *websocket.Conn
}

func NewClimate(conn *websocket.Conn) *Climate {
	return &Climate{
		conn: conn,
	}
}

/* Public API */

func (c Climate) SetFanMode(entityId string, fanMode string) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "climate"
	req.Service = "set_fan_mode"
	req.ServiceData = map[string]any{"fan_mode": fanMode}

	c.conn.WriteMessage(req)
}

func (c Climate) SetTemperature(entityId string, serviceData types.SetTemperatureRequest) {
	req := NewBaseServiceRequest(c.conn, entityId)
	req.Domain = "climate"
	req.Service = "set_temperature"
	req.ServiceData = serviceData.ToJSON()

	c.conn.WriteMessage(req)
}
