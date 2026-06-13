package gomeassistant

import (
	"saml.dev/gome-assistant/internal/services"
)

type Target = services.Target

func Entities(entityIDs ...string) Target {
	return services.Entities(entityIDs...)
}

func Areas(areaIDs ...string) Target {
	return services.Areas(areaIDs...)
}

func Devices(deviceIDs ...string) Target {
	return services.Devices(deviceIDs...)
}
