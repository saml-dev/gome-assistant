package services

/* Structs */

type Lock struct {
	api API
}

/* Public API */

// Lock a lock entity. Takes an entityID and an optional service_data,
// which must be serializable to a JSON object.
func (l Lock) Lock(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "lock",
		Service:     "lock",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return l.api.CallAndForget(req)
}

// Unlock a lock entity. Takes an entityID and an optional
// service_data, which must be serializable to a JSON object.
func (l Lock) Unlock(entityID string, serviceData ...any) error {
	req := BaseServiceRequest{
		Domain:      "lock",
		Service:     "unlock",
		ServiceData: optionalServiceData(serviceData...),
		Target:      Entity(entityID),
	}

	return l.api.CallAndForget(req)
}
