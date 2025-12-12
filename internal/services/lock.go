package services

/* Structs */

type Lock struct {
	api API
}

/* Public API */

// Lock a lock entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Lock) Lock(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "lock",
		Service: "lock",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return l.api.Call(req)
}

// Unlock a lock entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Lock) Unlock(entityID string, serviceData ...map[string]any) error {
	req := BaseServiceRequest{
		Domain:  "lock",
		Service: "unlock",
		Target:  Entity(entityID),
	}
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}
	return l.api.Call(req)
}
