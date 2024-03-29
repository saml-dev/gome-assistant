package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type ZWaveJS struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// ZWaveJS bulk_set_partial_config_parameters service.
func (zw ZWaveJS) BulkSetPartialConfigParam(entityId string, parameter int, value any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "zwave_js"
	req.Service = "bulk_set_partial_config_parameters"
	req.ServiceData = map[string]any{
		"parameter": parameter,
		"value":     value,
	}

	zw.conn.WriteMessage(req)
}
