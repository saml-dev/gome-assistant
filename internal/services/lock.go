package services

/* Structs */

type Lock struct {
	api API
}

/* Public API */

// Lock a lock entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Lock) Lock(entityID string, serviceData any) error {
	req := BaseServiceRequest{
		Domain:      "lock",
		Service:     "lock",
		Target:      Entity(entityID),
		ServiceData: serviceData,
	}

	return l.api.CallAndForget(req)
}

// Unlock a lock entity. Takes an entityID and an optional
// map that is translated into service_data.
func (l Lock) Unlock(entityID string, serviceData any) error {
	req := BaseServiceRequest{
		Domain:      "lock",
		Service:     "unlock",
		Target:      Entity(entityID),
		ServiceData: serviceData,
	}

	return l.api.CallAndForget(req)
}
