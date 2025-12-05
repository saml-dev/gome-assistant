package services

/* Structs */

type InputButton struct {
	api API
}

/* Public API */

func (ib InputButton) Press(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "input_button",
		Service: "press",
		Target:  Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputButton) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_button",
		Service: "reload",
		Target:  Entity(""),
	}
	return ib.api.Call(req)
}
