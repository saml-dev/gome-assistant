package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

type Notify struct {
	api API
}

// Notify sends a notification.
func (ha *Notify) Notify(
	ctx context.Context, serviceName string, reqData message.NotifyData,
) (any, error) {
	req := message.CallServiceData{
		Domain:      "notify",
		Service:     serviceName,
		ServiceData: reqData,
	}

	var result any
	if err := ha.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
