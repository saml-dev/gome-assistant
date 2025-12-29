package services

import (
	"saml.dev/gome-assistant/types"
)

type Notify struct {
	api API
}

// Notify sends a notification. Takes a types.NotifyRequest.
func (ha *Notify) Notify(reqData types.NotifyRequest) error {
	req := BaseServiceRequest{
		Domain:  "notify",
		Service: reqData.ServiceName,
	}
	serviceData := map[string]any{}
	serviceData["message"] = reqData.Message
	serviceData["title"] = reqData.Title
	if reqData.Data != nil {
		serviceData["data"] = reqData.Data
	}
	req.ServiceData = serviceData
	return ha.api.CallAndForget(req)
}
