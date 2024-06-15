package app

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-module/carbon"
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/websocket"
)

type State interface {
	Latitude() float64
	Longitude() float64
	AfterSunrise(...DurationString) bool
	BeforeSunrise(...DurationString) bool
	AfterSunset(...DurationString) bool
	BeforeSunset(...DurationString) bool
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
	EntityID    string              `json:"entity_id"`
	State       string              `json:"state"`
	Attributes  map[string]any      `json:"attributes"`
	LastChanged websocket.TimeStamp `json:"last_changed"`

	// The whole message, in JSON format:
	Raw websocket.RawMessage `json:"-"`
}

func newState(c *http.HttpClient, homeZoneEntityID string) (*StateImpl, error) {
	state := &StateImpl{httpClient: c}
	err := state.getLatLong(c, homeZoneEntityID)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (s *StateImpl) getLatLong(c *http.HttpClient, homeZoneEntityID string) error {
	resp, err := s.Get(homeZoneEntityID)
	if err != nil {
		return fmt.Errorf(
			"couldn't get latitude/longitude from home assistant entity '%s'. "+
				"Did you type it correctly? It should be a zone like 'zone.home'",
			homeZoneEntityID,
		)
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

func (s *StateImpl) Latitude() float64 {
	return s.latitude
}

func (s *StateImpl) Longitude() float64 {
	return s.longitude
}

func (s *StateImpl) Get(entityID string) (EntityState, error) {
	resp, err := s.httpClient.GetState(entityID)
	if err != nil {
		return EntityState{}, err
	}
	es := EntityState{}
	json.Unmarshal(resp, &es)
	es.Raw = resp
	return es, nil
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
