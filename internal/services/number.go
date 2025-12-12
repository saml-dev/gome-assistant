package services

type Number struct {
	api API
}

func (ib Number) SetValue(entityID string, value float32) error {
	req := BaseServiceRequest{
		Domain:      "number",
		Service:     "set_value",
		ServiceData: map[string]any{"value": value},
		Target:      Entity(entityID),
	}
	return ib.api.Call(req)
}

func (ib Number) MustSetValue(entityID string, value float32) {
	if err := ib.SetValue(entityID, value); err != nil {
		panic(err)
	}
}
