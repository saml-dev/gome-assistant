package services

/* Structs */

type InputButton struct {
	api API
}

/* Public API */

func (ib InputButton) Press(entityId string) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_button"
	req.Service = "press"

	return ib.api.WriteMessage(req)
}

func (ib InputButton) Reload() error {
	req := NewBaseServiceRequest("")
	req.Domain = "input_button"
	req.Service = "reload"
	return ib.api.WriteMessage(req)
}
