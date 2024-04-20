package services

import (
	"saml.dev/gome-assistant/internal/websocket"
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

type SetTemperatureRequest struct {
	Temperature    *float32
	TargetTempHigh *float32
	TargetTempLow  *float32
	HvacMode       string
}

func (r *SetTemperatureRequest) ToJSON() map[string]any {
	m := map[string]any{}
	if r.Temperature != nil {
		m["temperature"] = *r.Temperature
	}
	if r.TargetTempHigh != nil {
		m["target_temp_high"] = *r.TargetTempHigh
	}
	if r.TargetTempLow != nil {
		m["target_temp_low"] = *r.TargetTempLow
	}
	if r.HvacMode != "" {
		m["hvac_mode"] = r.HvacMode
	}
	return m
}

func (c Climate) SetTemperature(entityID string, setTemperatureRequest SetTemperatureRequest) {
	req := CallServiceRequest{
		Domain:  "climate",
		Service: "set_temperature",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: setTemperatureRequest.ToJSON(),
	}

	c.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
