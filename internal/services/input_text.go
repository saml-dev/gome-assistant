package services

/* Structs */

type InputText struct {
	api API
}

/* Public API */

func (ib InputText) Set(entityId string, value string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_text"
	req.Service = "set_value"
	req.ServiceData = map[string]any{
		"value": value,
	}

	return ib.api.WriteMessage(req)
}

func (ib InputText) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "input_text"
	req.Service = "reload"
	return ib.api.WriteMessage(req)
}
