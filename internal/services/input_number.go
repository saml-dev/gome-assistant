package services

/* Structs */

type InputNumber struct {
	api API
}

/* Public API */

func (ib InputNumber) Set(entityID string, value float32) error {
	req := BaseServiceRequest{
		Domain:      "input_number",
		Service:     "set_value",
		ServiceData: map[string]any{"value": value},
		Target:      Entity(entityID),
	}
	return ib.api.Call(req)
}

func (ib InputNumber) Increment(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "increment",
		Target:  Entity(entityID),
	}
	return ib.api.Call(req)
}

func (ib InputNumber) Decrement(entityID string) error {
	req := BaseServiceRequest{
		Domain:  "input_number",
		Service: "decrement",
		Target:  Entity(entityID),
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
