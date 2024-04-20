package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

type Notify struct {
	conn *websocket.Conn
}

func NewNotify(conn *websocket.Conn) *Notify {
	return &Notify{
		conn: conn,
	}
}

type NotifyRequest struct {
	// Which notify service to call, such as mobile_app_sams_iphone
	ServiceName string
	Message     string
	Title       string
	Data        map[string]any
}

// Send a notification.
func (ha *Notify) Notify(reqData NotifyRequest) {
	serviceData := map[string]any{
		"message": reqData.Message,
		"title":   reqData.Title,
	}
	if reqData.Data != nil {
		serviceData["data"] = reqData.Data
	}

	req := CallServiceRequest{
		Domain:      "notify",
		Service:     reqData.ServiceName,
		ServiceData: serviceData,
	}

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
