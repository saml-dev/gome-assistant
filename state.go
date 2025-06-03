package gomeassistant

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-module/carbon"

	"saml.dev/gome-assistant/internal/http"
)

type State interface {
	AfterSunrise(...DurationString) bool
	BeforeSunrise(...DurationString) bool
	AfterSunset(...DurationString) bool
	BeforeSunset(...DurationString) bool
	ListEntities() ([]EntityState, error)
	Get(entityId string) (EntityState, error)
	Equals(entityId, state string) (bool, error)
}

// State is used to retrieve state from Home Assistant.
type StateImpl struct {
	httpClient *http.HttpClient
	latitude   float64
	longitude  float64
}

type EntityState struct {
	EntityID    string         `json:"entity_id"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastChanged time.Time      `json:"last_changed"`
}

func newState(c *http.HttpClient, homeZoneEntityId string) (*StateImpl, error) {
	state := &StateImpl{httpClient: c}

	// Ensure the zone exists and has required attributes
	entity, err := state.Get(homeZoneEntityId)
	if err != nil {
		return nil, fmt.Errorf("home zone entity '%s' not found: %w", homeZoneEntityId, err)
	}

	// Ensure it's a zone entity
	if !strings.HasPrefix(homeZoneEntityId, "zone.") {
		return nil, fmt.Errorf("entity '%s' is not a zone entity (must start with zone.)", homeZoneEntityId)
	}

	// Verify and extract latitude and longitude
	if entity.Attributes == nil {
		return nil, fmt.Errorf("home zone entity '%s' has no attributes", homeZoneEntityId)
	}

	if lat, ok := entity.Attributes["latitude"].(float64); ok {
		state.latitude = lat
	} else {
		return nil, fmt.Errorf("home zone entity '%s' missing valid latitude attribute", homeZoneEntityId)
	}

	if long, ok := entity.Attributes["longitude"].(float64); ok {
		state.longitude = long
	} else {
		return nil, fmt.Errorf("home zone entity '%s' missing valid longitude attribute", homeZoneEntityId)
	}

	return state, nil
}

func (s *StateImpl) Get(entityId string) (EntityState, error) {
	resp, err := s.httpClient.GetState(entityId)
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

func (s *StateImpl) Equals(entityId string, expectedState string) (bool, error) {
	currentState, err := s.Get(entityId)
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
