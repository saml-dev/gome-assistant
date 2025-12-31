package services

import "context"

type HomeAssistant struct {
	api API
}

// TurnOn a Home Assistant entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (ha *HomeAssistant) TurnOn(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "homeassistant",
		Service:     "turn_on",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	var result any
	if err := ha.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (ha *HomeAssistant) Toggle(
	ctx context.Context, entityID string, serviceData ...any,
) (any, error) {
	req := BaseServiceRequest{
		Domain:      "homeassistant",
		Service:     "toggle",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	var result any
	if err := ha.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (ha *HomeAssistant) TurnOff(
	ctx context.Context, entityID string,
) (any, error) {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_off",
		Target:  Entity(entityID),
	}

	var result any
	if err := ha.api.Call(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
