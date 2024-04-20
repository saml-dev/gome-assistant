package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Lock struct {
	conn *websocket.Conn
}

func NewLock(conn *websocket.Conn) *Lock {
	return &Lock{
		conn: conn,
	}
}

/* Public API */

// Lock a lock entity.
func (l Lock) Lock(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "lock",
		Service: "lock",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

// Unlock a lock entity.
func (l Lock) Unlock(entityID string, serviceData map[string]any) {
	req := CallServiceRequest{
		Domain:  "lock",
		Service: "unlock",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	l.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
