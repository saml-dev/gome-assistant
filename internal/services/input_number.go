package services

/* Structs */

type InputNumber struct {
	api API
}

/* Public API */

func (ib InputNumber) Set(entityId string, value float32) error {
	req := BaseServiceRequest{
		Domain:      "input_number",
		Service:     "set_value",
		ServiceData: map[string]any{"value": value},
		Target:      Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputNumber) Increment(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "increment",
		Target:  Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputNumber) Decrement(entityId string) error {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "decrement",
		Target:  Entity(entityId),
	}
	return ib.api.Call(req)
}

func (ib InputNumber) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "reload",
	}
	return ib.api.Call(req)
}
