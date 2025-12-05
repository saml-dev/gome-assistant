package services

/* Structs */

type InputText struct {
	api API
}

/* Public API */

func (ib InputText) Set(entityId string, value string) error {
	req := BaseServiceRequest{
		Domain:  "input_text",
		Service: "set_value",
		ServiceData: map[string]any{
			"value": value,
		},
		Target: Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputText) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_text",
		Service: "reload",
	}
	return ib.api.Call(req)
}
