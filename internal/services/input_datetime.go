package services

import (
	"fmt"
	"time"
)

/* Structs */

type InputDatetime struct {
	api API
}

/* Public API */

func (ib InputDatetime) Set(entityID string, value time.Time) error {
	req := BaseServiceRequest{
		Domain:  "input_datetime",
		Service: "set_datetime",
		ServiceData: map[string]any{
			"timestamp": fmt.Sprint(value.Unix()),
		},
		Target: Entity(entityID),
	}
	return ib.api.Call(req)
}

func (ib InputDatetime) Reload() error {
	req := BaseServiceRequest{
		Domain:  "input_datetime",
		Service: "reload",
	}
	return ib.api.Call(req)
}
