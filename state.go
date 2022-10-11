package gomeassistant

import "github.com/saml-dev/gome-assistant/internal/http"

// State is used to retrieve state from Home Assistant.
type State struct {
	httpClient *http.HttpClient
}

func NewState(c *http.HttpClient) *State {
	return &State{httpClient: c}
}

func (s *State) Get(entityId string) (string, error) {
	resp, err := s.httpClient.GetState(entityId)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
