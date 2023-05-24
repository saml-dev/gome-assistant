package gomeassistant

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-module/carbon"
	"saml.dev/gome-assistant/internal/http"
)

// State is used to retrieve state from Home Assistant.
type State struct {
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

func newState(c *http.HttpClient, homeZoneEntityId string) (*State, error) {
	state := &State{httpClient: c}
	err := state.getLatLong(c, homeZoneEntityId)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (s *State) getLatLong(c *http.HttpClient, homeZoneEntityId string) error {
	resp, err := s.Get(homeZoneEntityId)
	if err != nil {
		return errors.New(fmt.Sprintf("couldn't get latitude/longitude from home assistant entity '%s'. Did you type it correctly? It should be a zone like 'zone.home'.\n", homeZoneEntityId))
	}

	if resp.Attributes["latitude"] != nil {
		s.latitude = resp.Attributes["latitude"].(float64)
	} else {
		return errors.New("server returned nil latitude")
	}

	if resp.Attributes["longitude"] != nil {
		s.longitude = resp.Attributes["longitude"].(float64)
	} else {
		return errors.New("server returned nil longitude")
	}

	return nil
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
	sunrise := getSunriseSunset(s /* sunrise = */, true, carbon.Now(), offset...)
	return carbon.Now().Lt(sunrise)
}

func (s *State) AfterSunrise(offset ...DurationString) bool {
	return !s.BeforeSunrise(offset...)
}

func (s *State) BeforeSunset(offset ...DurationString) bool {
	sunset := getSunriseSunset(s /* sunrise = */, false, carbon.Now(), offset...)
	return carbon.Now().Lt(sunset)
}

func (s *State) AfterSunset(offset ...DurationString) bool {
	return !s.BeforeSunset(offset...)
}
