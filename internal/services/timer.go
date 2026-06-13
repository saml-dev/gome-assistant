package services

import "context"

/* Structs */

type Timer struct {
	api API
}

/* Public API */

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Start(
	ctx context.Context, target Target, duration string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "start",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: target,
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerstart
func (t Timer) Change(
	ctx context.Context, target Target, duration string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "change",
		ServiceData: map[string]any{
			"duration": duration,
		},
		Target: target,
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerpause
func (t Timer) Pause(
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "pause",
		Target:  target,
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timercancel
func (t Timer) Cancel(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
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
	ctx context.Context, target Target,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "finish",
		Target:  target,
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// See https://www.home-assistant.io/integrations/timer/#action-timerreload
func (t Timer) Reload(ctx context.Context) (any, error) {
	req := BaseServiceRequest{
		Domain:  "timer",
		Service: "reload",
	}

	var result any
	if err := t.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
