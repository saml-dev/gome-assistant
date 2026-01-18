package gomeassistant

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-module/carbon"

	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/message"
)

type State interface {
	AfterSunrise(...DurationString) bool
	BeforeSunrise(...DurationString) bool
	AfterSunset(...DurationString) bool
	BeforeSunset(...DurationString) bool
	ListEntities() ([]EntityState, error)
	Get(entityID string) (EntityState, error)
	Equals(entityID, state string) (bool, error)
}

// State is used to retrieve state from Home Assistant.
type StateImpl struct {
	httpClient *http.HttpClient
	latitude   float64
	longitude  float64
}

type EntityState struct {
	EntityID    string            `json:"entity_id"`
	State       string            `json:"state"`
	Attributes  map[string]any    `json:"attributes"`
	LastChanged message.TimeStamp `json:"last_changed"`
}

func newState(c *http.HttpClient, homeZoneEntityID string) (*StateImpl, error) {
	state := &StateImpl{httpClient: c}

	// Ensure the zone exists and has required attributes
	entity, err := state.Get(homeZoneEntityID)
	if err != nil {
		return nil, fmt.Errorf("home zone entity '%s' not found: %w", homeZoneEntityID, err)
	}

	// Ensure it's a zone entity
	if !strings.HasPrefix(homeZoneEntityID, "zone.") {
		return nil, fmt.Errorf("entity '%s' is not a zone entity (must start with zone.)", homeZoneEntityID)
	}

	// Verify and extract latitude and longitude
	if entity.Attributes == nil {
		return nil, fmt.Errorf("home zone entity '%s' has no attributes", homeZoneEntityID)
	}

	if lat, ok := entity.Attributes["latitude"].(float64); ok {
		state.latitude = lat
	} else {
		return nil, fmt.Errorf("home zone entity '%s' missing valid latitude attribute", homeZoneEntityID)
	}

	if long, ok := entity.Attributes["longitude"].(float64); ok {
		state.longitude = long
	} else {
		return nil, fmt.Errorf("home zone entity '%s' missing valid longitude attribute", homeZoneEntityID)
	}

	return state, nil
}

func (s *StateImpl) Get(entityID string) (EntityState, error) {
	resp, err := s.httpClient.GetState(entityID)
	if err != nil {
		return EntityState{}, err
	}
	es := EntityState{}
	err = json.Unmarshal(resp, &es)
	return es, err
}

// ListEntities returns a list of all entities in Home Assistant.
// see rest documentation for more details: https://developers.home-assistant.io/docs/api/rest/#actions
func (s *StateImpl) ListEntities() ([]EntityState, error) {
	resp, err := s.httpClient.States()
	if err != nil {
		return nil, err
	}
	es := []EntityState{}
	err = json.Unmarshal(resp, &es)
	return es, err
}

func (s *StateImpl) Equals(entityID string, expectedState string) (bool, error) {
	currentState, err := s.Get(entityID)
	if err != nil {
		return false, err
	}
	return currentState.State == expectedState, nil
}

func (s *StateImpl) BeforeSunrise(offset ...DurationString) bool {
	sunrise := getSunriseSunset(s /* sunrise = */, true, carbon.Now(), offset...)
	return carbon.Now().Lt(sunrise)
}

func (s *StateImpl) AfterSunrise(offset ...DurationString) bool {
	return !s.BeforeSunrise(offset...)
}

func (s *StateImpl) BeforeSunset(offset ...DurationString) bool {
	sunset := getSunriseSunset(s /* sunrise = */, false, carbon.Now(), offset...)
	return carbon.Now().Lt(sunset)
}

func (s *StateImpl) AfterSunset(offset ...DurationString) bool {
	return !s.BeforeSunset(offset...)
}
