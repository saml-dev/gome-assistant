package main

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang-cz/devslog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	gaapp "saml.dev/gome-assistant/app"
)

type (
	MySuite struct {
		suite.Suite
		app      *gaapp.App
		config   *Config
		suiteCtx map[string]any
	}

	Config struct {
		Hass struct {
			HAAuthToken      string `yaml:"token"`
			IpAddress        string `yaml:"address"`
			Port             string `yaml:"port"`
			HomeZoneEntityID string `yaml:"zone"`
		}
		Entities struct {
			LightEntityID string `yaml:"light_entity_id"`
		}
	}
)

func setupLogging() {
	opts := &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
		NewLineAfterLog: true,
	}
	slog.SetDefault(slog.New(devslog.NewHandler(os.Stdout, opts)))
}

func (s *MySuite) SetupSuite(ctx context.Context) {
	setupLogging()
	slog.Debug("Setting up test suite...")
	s.suiteCtx = make(map[string]any)

	configFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		slog.Error("Error reading config file", err)
	}
	s.config = &Config{}
	// either env var or config file can be used to set HA auth. token
	s.config.Hass.HAAuthToken = os.Getenv("HA_AUTH_TOKEN")
	if err := yaml.Unmarshal(configFile, s.config); err != nil {
		slog.Error("Error unmarshalling config file", err)
	}

	s.app, err = gaapp.NewApp(
		ctx,
		gaapp.NewAppRequest{
			HAAuthToken:      s.config.Hass.HAAuthToken,
			IpAddress:        s.config.Hass.IpAddress,
			HomeZoneEntityID: s.config.Hass.HomeZoneEntityID,
		},
	)
	if err != nil {
		slog.Error("Failed to createw new app", err)
		s.T().FailNow()
	}

	// Register all automations
	entityID := s.config.Entities.LightEntityID
	if entityID != "" {
		s.suiteCtx["entityCallbackInvoked"] = false
		etl := gaapp.NewEntityListener().EntityIDs(entityID).Call(s.entityCallback).Build()
		s.app.RegisterEntityListeners(etl)
	}

	s.suiteCtx["dailyScheduleCallbackInvoked"] = false
	runTime := time.Now().Add(1 * time.Minute).Format("15:04")
	dailySchedule := gaapp.NewDailySchedule().Call(s.dailyScheduleCallback).At(runTime).Build()
	s.app.RegisterSchedules(dailySchedule)

	// start GA app
	go s.app.Start(ctx)
}

func (s *MySuite) TearDownSuite() {
	if s.app != nil {
		s.app.Close()
		s.app = nil
	}
}

// Basic test of light toggle service and entity listener
func (s *MySuite) TestLightService() {
	entityID := s.config.Entities.LightEntityID

	if entityID != "" {
		initState := getEntityState(s, entityID)
		s.app.GetService().Light.Toggle(entityID, nil)

		assert.EventuallyWithT(
			s.T(),
			func(c *assert.CollectT) {
				newState := getEntityState(s, entityID)
				assert.NotEqual(c, initState, newState)
				assert.True(c, s.suiteCtx["entityCallbackInvoked"].(bool))
			},
			10*time.Second, 1*time.Second,
			"State of light entity did not change or callback was not invoked",
		)
	} else {
		s.T().Skip("No light entity id provided")
	}
}

// Basic test of daily schedule and callback
func (s *MySuite) TestSchedule() {
	assert.EventuallyWithT(s.T(), func(c *assert.CollectT) {
		assert.True(c, s.suiteCtx["dailyScheduleCallbackInvoked"].(bool))
	}, 2*time.Minute, 1*time.Second, "Daily schedule callback was not invoked")
}

// Capture event after light entity state has changed
func (s *MySuite) entityCallback(se *gaapp.Service, st gaapp.State, e gaapp.EntityData) {
	slog.Info(
		"Entity callback called.",
		"entity id", e.TriggerEntityID,
		"from state", e.FromState,
		"to state", e.ToState,
	)
	s.suiteCtx["entityCallbackInvoked"] = true
}

// Capture planned daily schedule
func (s *MySuite) dailyScheduleCallback(se *gaapp.Service, st gaapp.State) {
	slog.Info("Daily schedule callback called.")
	s.suiteCtx["dailyScheduleCallbackInvoked"] = true
}

func getEntityState(s *MySuite, entityID string) string {
	state, err := s.app.GetState().Get(entityID)
	if err != nil {
		slog.Error("Error getting entity state", err)
		s.T().FailNow()
	}
	slog.Info("State of entity", "state", state.State)
	return state.State
}

// Run the test suite
func TestMySuite(t *testing.T) {
	suite.Run(t, new(MySuite))
}
