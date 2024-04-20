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

func (c Climate) SetFanMode(entityID string, fanMode string) {
	req := CallServiceRequest{
		Domain:  "climate",
		Service: "set_fan_mode",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: map[string]any{"fan_mode": fanMode},
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (c Climate) SetTemperature(entityID string, serviceData types.SetTemperatureRequest) {
	req := CallServiceRequest{
		Domain:  "climate",
		Service: "set_temperature",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData.ToJSON(),
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
