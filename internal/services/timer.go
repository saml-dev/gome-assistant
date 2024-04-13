package services

import (
	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Timer struct {
	conn *ws.WebsocketConn
}

/* Public API */

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Start(entityId string, duration string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "timer"
	req.Service = "start"
	req.ServiceData = map[string]any{
		"duration": duration,
	}

	return t.conn.WriteMessage(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Change(entityId string, duration string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "timer"
	req.Service = "change"
	req.ServiceData = map[string]any{
		"duration": duration,
	}

	return t.conn.WriteMessage(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerpause
func (t Timer) Pause(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "timer"
	req.Service = "pause"
	return t.conn.WriteMessage(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timercancel
func (t Timer) Cancel() error {
	req := NewBaseServiceRequest("")
	req.Domain = "timer"
	req.Service = "cancel"
	return t.conn.WriteMessage(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerfinish
func (t Timer) Finish(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "timer"
	req.Service = "finish"
	return t.conn.WriteMessage(req)
}

// See https://www.home-assistant.io/integrations/timer/#action-timerreload
func (t Timer) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "timer"
	req.Service = "reload"
	return t.conn.WriteMessage(req)
}
