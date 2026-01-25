package services

import (
	"context"

	"saml.dev/gome-assistant/message"
)

/* Structs */

type Timer struct {
	api API
}

/* Public API */

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Start(
	ctx context.Context, entityID string, duration string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "timer",
		Service: "start",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: message.Entity(entityID),
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Change(
	ctx context.Context, entityID string, duration string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "timer",
		Service: "change",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: message.Entity(entityID),
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerpause
func (t Timer) Pause(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "timer",
		Service: "pause",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timercancel
func (t Timer) Cancel(ctx context.Context) (any, error) {
	req := message.CallServiceData{
		Domain:  "timer",
		Service: "cancel",
	}
	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerfinish
func (t Timer) Finish(
	ctx context.Context, entityID string,
) (any, error) {
	req := message.CallServiceData{
		Domain:  "timer",
		Service: "finish",
		Target:  message.Entity(entityID),
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerreload
func (t Timer) Reload(ctx context.Context) (any, error) {
	req := message.CallServiceData{
		Domain:  "timer",
		Service: "reload",
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
