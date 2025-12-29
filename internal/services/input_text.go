package services

/* Structs */

type InputText struct {
	api API
}

/* Public API */

func (ib InputText) Set(entityID string, value string) error {
	req := BaseServiceRequest{
		Domain:  "input_text",
		Service: "set_value",
		ServiceData: map[string]any{
			"value": value,
		},
		Target: Entity(entityID),
	}
	return ib.api.CallAndForget(req)
}

func (ib InputText) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_text",
		Service: "reload",
	}
	return ib.api.CallAndForget(req)
}
