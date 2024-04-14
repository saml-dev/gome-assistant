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

// Send a notification. Takes a types.NotifyRequest.
func (ha *Notify) Notify(reqData types.NotifyRequest) {
	req := NewBaseServiceRequest(ha.conn, "")
	req.Domain = "notify"
	req.Service = reqData.ServiceName

	serviceData := map[string]any{}
	serviceData["message"] = reqData.Message
	serviceData["title"] = reqData.Title
	if reqData.Data != nil {
		serviceData["data"] = reqData.Data
	}

	req.ServiceData = serviceData

	ha.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
