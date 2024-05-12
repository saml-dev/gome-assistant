package ga

import "fmt"

// Target represents the target of the service call, if applicable.
type Target struct {
	EntityID string `json:"entity_id,omitempty"`
	DeviceID string `json:"device_id,omitempty"`
}

func EntityTarget(entityID string) Target {
	return Target{
		EntityID: entityID,
	}
}

func DeviceTarget(deviceID string) Target {
	return Target{
		DeviceID: deviceID,
	}
}

func (t Target) String() string {
	switch {
	case t.EntityID != "":
		return fmt.Sprintf("entity %s", t.EntityID)
	case t.DeviceID != "":
		return fmt.Sprintf("device %s", t.DeviceID)
	default:
		return "unset target"
	}
}
