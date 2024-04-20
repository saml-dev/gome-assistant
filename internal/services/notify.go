package services

import (
	"saml.dev/gome-assistant/internal/websocket"
	"saml.dev/gome-assistant/types"
)

type Notify struct {
	conn *websocket.Conn
}

func NewNotify(conn *websocket.Conn) *Notify {
	return &Notify{
		conn: conn,
	}
}

// Send a notification.
func (ha *Notify) Notify(reqData types.NotifyRequest) {
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
