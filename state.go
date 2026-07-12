package gomeassistant

import (
	"context"
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
	Get(entityID string) (EntityState, error)
	GetStates(entityIDs []string) ([]string, error)
	Equals(entityID, state string) (bool, error)
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

func newState(c *http.HttpClient) *StateImpl {
	return &StateImpl{httpClient: c}
}

// loadHomeZone loads and validates the home zone coordinates for sunrise and
// sunset calculations. Call this during application startup, not construction.
func (s *StateImpl) loadHomeZone(ctx context.Context, homeZoneEntityID string) error {
	if !strings.HasPrefix(homeZoneEntityID, "zone.") {
		return fmt.Errorf("entity '%s' is not a zone entity (must start with zone.)", homeZoneEntityID)
	}

	entity, err := s.getWithContext(ctx, homeZoneEntityID)
	if err != nil {
		return fmt.Errorf("home zone entity '%s' not found: %w", homeZoneEntityID, err)
	}

	if entity.Attributes == nil {
		return fmt.Errorf("home zone entity '%s' has no attributes", homeZoneEntityID)
	}

	lat, ok := entity.Attributes["latitude"].(float64)
	if !ok {
		return fmt.Errorf("home zone entity '%s' missing valid latitude attribute", homeZoneEntityID)
	}

	longitude, ok := entity.Attributes["longitude"].(float64)
	if !ok {
		return fmt.Errorf("home zone entity '%s' missing valid longitude attribute", homeZoneEntityID)
	}

	s.latitude = lat
	s.longitude = longitude
	return nil
}

func (s *StateImpl) Get(entityID string) (EntityState, error) {
	return s.getWithContext(context.Background(), entityID)
}

func (s *StateImpl) getWithContext(ctx context.Context, entityID string) (EntityState, error) {
	resp, err := s.httpClient.GetStateWithContext(ctx, entityID)
	if err != nil {
		return EntityState{}, err
	}
	es := EntityState{}
	err = json.Unmarshal(resp, &es)
	return es, err
}

// GetStates returns the states for each entity ID in the same order as
// entityIDs. All values are strings; callers must convert returned values to
// numbers or other types when applicable.
func (s *StateImpl) GetStates(entityIDs []string) ([]string, error) {
	if len(entityIDs) == 0 {
		return []string{}, nil
	}

	resp, err := s.httpClient.GetStates(entityIDs)
	if err != nil {
		return nil, err
	}
	states := []string{}
	err = json.Unmarshal(resp, &states)
	return states, err
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
