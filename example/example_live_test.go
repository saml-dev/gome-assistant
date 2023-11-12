package main

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	ga "saml.dev/gome-assistant"
)

type (
	MySuite struct {
		suite.Suite
		app    *ga.App
		config *Config
	}

	Config struct {
		App struct {
			HAAuthToken      string `yaml:"ha_auth_token"`
			IpAddress        string `yaml:"ip_address"`
			Port             string `yaml:"port"`
			HomeZoneEntityId string `yaml:"home_zone_entity_id"`
		}
		Entities struct {
			LightEntityId string `yaml:"light_entity_id"`
		}
	}
)

func (s *MySuite) SetupTest() {
	configFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		slog.Error("Error reading config file", err)
	}
	s.config = &Config{}
	// either env var or config file can be used to set HA auth. token
	s.config.App.HAAuthToken = os.Getenv("HA_AUTH_TOKEN")
	if err := yaml.Unmarshal(configFile, s.config); err != nil {
		slog.Error("Error unmarshalling config file:", err)
	}

	s.app, err = ga.NewApp(ga.NewAppRequest{
		HAAuthToken:      s.config.App.HAAuthToken,
		IpAddress:        s.config.App.IpAddress,
		HomeZoneEntityId: s.config.App.HomeZoneEntityId,
	})
	if err != nil {
		slog.Error("Failed to createw new app:", err)
		s.T().FailNow()
	}
}

func (s *MySuite) TearDownSuite() {
	if s.app != nil {
		s.app.Cleanup()
	}
}

// Basic test of ga app creation and light toggle service
func (s *MySuite) TestLightService() {
	entityId := s.config.Entities.LightEntityId

	initialState, err := s.app.GetState().Get(entityId)
	if err != nil {
		slog.Error("Error getting entity state:", err)
	}
	slog.Info("Initial state of entity:", "state", initialState.State)

	s.app.GetService().Light.Toggle(entityId)

	time.Sleep(1 * time.Second) // wait for state to update

	newState, err := s.app.GetState().Get(entityId)
	if err != nil {
		slog.Error("Error getting entity state:", err)
	}
	slog.Info("New state of entity:", "state", newState.State)
	assert.Equal(s.T(), initialState.State, newState.State)
}

// Run the test suite
func TestMySuite(t *testing.T) {
	suite.Run(t, new(MySuite))
}
