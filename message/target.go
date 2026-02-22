package message

type Target struct {
	EntityID string `json:"entity_id,omitempty"`
}

func Entity(entityID string) Target {
	return Target{
		EntityID: entityID,
	}
}
