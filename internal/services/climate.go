package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
	"saml.dev/gome-assistant/types"
)

/* Structs */

type Climate struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

func (c Climate) SetFanMode(entityId string, fanMode string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "climate"
	req.Service = "set_fan_mode"
	req.ServiceData = map[string]any{"fan_mode": fanMode}

	c.conn.WriteMessage(c.ctx, req)
}

func (c Climate) SetTemperature(entityId string, serviceData types.SetTemperatureRequest) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "climate"
	req.Service = "set_temperature"
	req.ServiceData = serviceData.ToJSON()

	c.conn.WriteMessage(c.ctx, req)
}
