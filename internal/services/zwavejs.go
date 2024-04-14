package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type ZWaveJS struct {
	conn *websocket.Conn
}

func NewZWaveJS(conn *websocket.Conn) *ZWaveJS {
	return &ZWaveJS{
		conn: conn,
	}
}

/* Public API */

// ZWaveJS bulk_set_partial_config_parameters service.
func (zw ZWaveJS) BulkSetPartialConfigParam(entityID string, parameter int, value any) {
	req := CallServiceRequest{}
	req.Domain = "zwave_js"
	req.Service = "bulk_set_partial_config_parameters"
	req.Target.EntityID = entityID
	req.ServiceData = map[string]any{
		"parameter": parameter,
		"value":     value,
	}

	zw.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
