package services

import (
	"context"

	"saml.dev/gome-assistant/message"
	"saml.dev/gome-assistant/types"
)

type Notify struct {
	api API
}

// Notify sends a notification. Takes a types.NotifyRequest.
func (ha *Notify) Notify(
	ctx context.Context, reqData types.NotifyRequest,
) (any, error) {
	req := message.CallServiceData{
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

	var result any
	if err := ha.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
