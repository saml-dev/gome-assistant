package services

import (
	"context"

	"github.com/gorilla/websocket"
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Lock struct {
	conn *websocket.Conn
	ctx  context.Context
}

/* Public API */

// Lock a lock entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Lock) Lock(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "lock"
	req.Service = "lock"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, l.conn, l.ctx)
}

// Unlock a lock entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Lock) Unlock(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "lock"
	req.Service = "unlock"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	ws.WriteMessage(req, l.conn, l.ctx)
}
