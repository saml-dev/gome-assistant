package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
	"saml.dev/gome-assistant/types"
)

/* Structs */

type Climate struct {
	conn *ws.WebsocketConn
}

/* Public API */

func (c Climate) SetFanMode(entityId string, fanMode string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "climate"
	req.Service = "set_fan_mode"
	req.ServiceData = map[string]any{"fan_mode": fanMode}

	return c.conn.WriteMessage(req)
}

func (c Climate) SetTemperature(entityId string, serviceData types.SetTemperatureRequest) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "climate"
	req.Service = "set_temperature"
	req.ServiceData = serviceData.ToJSON()

	return c.conn.WriteMessage(req)
}
