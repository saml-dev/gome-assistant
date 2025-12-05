package services

type Number struct {
	api API
}

func (ib Number) SetValue(entityId string, value float32) error {
	req := BaseServiceRequest{
		Domain:      "number",
		Service:     "set_value",
		ServiceData: map[string]any{"value": value},
		Target:      Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib Number) MustSetValue(entityId string, value float32) {
	if err := ib.SetValue(entityId, value); err != nil {
		panic(err)
	}
}
