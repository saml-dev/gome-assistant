package services

import (
	"context"

	"saml.dev/gome-assistant/websocket"
)

type Notify struct {
	service Service
}

func NewNotify(service Service) *Notify {
	return &Notify{
		service: service,
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
func (ha *Notify) Notify(reqData NotifyRequest) (websocket.Message, error) {
	ctx := context.TODO()
	serviceData := map[string]any{
		"message": reqData.Message,
		"title":   reqData.Title,
	}
	if reqData.Data != nil {
		serviceData["data"] = reqData.Data
	}

	return ha.service.CallService(
		ctx, "notify", reqData.ServiceName,
		serviceData, Target{},
	)
}
