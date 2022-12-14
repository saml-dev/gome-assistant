package gomeassistant

import (
	"encoding/json"
	"time"

	"saml.dev/gome-assistant/internal/http"
)

// State is used to retrieve state from Home Assistant.
type State struct {
	httpClient *http.HttpClient
}

type EntityState struct {
	EntityID    string         `json:"entity_id"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
	LastChanged time.Time      `json:"last_changed"`
}

func newState(c *http.HttpClient) *State {
	return &State{httpClient: c}
}

func (s *State) Get(entityId string) (EntityState, error) {
	resp, err := s.httpClient.GetState(entityId)
	if err != nil {
		return EntityState{}, err
	}
	es := EntityState{}
	json.Unmarshal(resp, &es)
	return es, nil
}

func (s *State) Equals(entityId string, expectedState string) (bool, error) {
	currentState, err := s.Get(entityId)
	if err != nil {
		return false, err
	}
	return currentState.State == expectedState, nil
}

func (s *State) BeforeSunrise(offset ...DurationString) bool {
	sunrise := getSunriseSunsetFromState(s, true, offset...)
	return sunrise.IsToday()
}

func (s *State) AfterSunrise(offset ...DurationString) bool {
	return !s.BeforeSunrise(offset...)
}

func (s *State) BeforeSunset(offset ...DurationString) bool {
	sunset := getSunriseSunsetFromState(s, false, offset...)
	return sunset.IsToday()
}

func (s *State) AfterSunset(offset ...DurationString) bool {
	return !s.BeforeSunset(offset...)
}
