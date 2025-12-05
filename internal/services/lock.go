package services

/* Structs */

type Lock struct {
	api API
}

/* Public API */

// Lock a lock entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Lock) Lock(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "lock"
	req.Service = "lock"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return l.api.WriteMessage(req)
}

// Unlock a lock entity. Takes an entityId and an optional
// map that is translated into service_data.
func (l Lock) Unlock(entityId string, serviceData ...map[string]any) error {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "lock"
	req.Service = "unlock"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	return l.api.WriteMessage(req)
}
