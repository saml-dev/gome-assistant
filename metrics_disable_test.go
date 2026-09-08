package gomeassistant

import (
	"context"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
)

func TestMetricsDisabledAutomationsStillRunWithoutTracking(t *testing.T) {
	app, recorder := newDirectMetricsTestApp(t)
	type callbackObjects struct {
		kind    string
		service *Service
		state   State
	}
	callbacks := make(chan callbackObjects, 4)
	recordCallback := func(kind string, service *Service, state State) {
		callbacks <- callbackObjects{kind: kind, service: service, state: state}
	}

	daily := NewDailySchedule().
		Call(func(service *Service, state State) { recordCallback("daily", service, state) }).
		At("12:00").
		DisableMetrics().
		Build()
	interval := NewInterval().
		Call(func(service *Service, state State) { recordCallback("interval", service, state) }).
		Every("5s").
		DisableMetrics().
		Build()
	entity := NewEntityListener().
		EntityIDs("light.test").
		Call(func(service *Service, state State, _ EntityData) { recordCallback("entity", service, state) }).
		DisableMetrics().
		Build()
	event := NewEventListener().
		EventTypes("test").
		Call(func(service *Service, state State, _ EventData) { recordCallback("event", service, state) }).
		DisableMetrics().
		Build()

	daily.maybeRunCallback(app)
	interval.maybeRunCallback(app)
	entity.maybeCall(app,
		EntityData{TriggerEntityID: "light.test", FromState: "off", ToState: "on"},
		stateData{EntityID: "light.test", OldState: msgState{State: "off"}, NewState: msgState{State: "on"}},
	)
	event.maybeCall(app, EventData{Type: "test"})

	seen := make(map[string]bool, 4)
	for range 4 {
		select {
		case result := <-callbacks:
			seen[result.kind] = true
			require.Same(t, app.service, result.service)
			require.Same(t, app.state, result.state)
		case <-time.After(time.Second):
			t.Fatal("metrics-disabled automation callback was not invoked")
		}
	}
	app.workers.Wait()
	require.Equal(t, map[string]bool{"daily": true, "interval": true, "entity": true, "event": true}, seen)

	events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{Limit: 100})
	require.NoError(t, err)
	require.Empty(t, events)
}

func TestMetricsDisabledDelayedListenerDoesNotAffectWaitingStatus(t *testing.T) {
	app, recorder := newDirectMetricsTestApp(t)
	called := make(chan struct{})
	listener := NewEntityListener().
		EntityIDs("light.test").
		Call(func(*Service, State, EntityData) { close(called) }).
		ToState("on").
		Duration("10ms").
		DisableMetrics().
		Build()

	listener.maybeCall(app,
		EntityData{TriggerEntityID: "light.test", FromState: "off", ToState: "on"},
		stateData{EntityID: "light.test", OldState: msgState{State: "off"}, NewState: msgState{State: "on"}},
	)
	require.Zero(t, app.snapshotRuntimeStatus().Waiting)
	requireCallback(t, called)
	app.workers.Wait()
	require.Zero(t, app.snapshotRuntimeStatus().Waiting)

	events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{Limit: 100})
	require.NoError(t, err)
	require.Empty(t, events)
}

func TestMetricsDisabledRunOnStartupRunsOncePerSessionWithoutTracking(t *testing.T) {
	server, _, _ := newFakeHA(t, 200)
	app := fakeApp(t, server.URL)
	var callbackCount atomic.Int32
	app.RegisterEntityListeners(NewEntityListener().
		EntityIDs("light.first", "light.second").
		Call(func(*Service, State, EntityData) { callbackCount.Add(1) }).
		RunOnStartup().
		Key("disabled-startup").
		DisableMetrics().
		Build())
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	runReadyMetricsSession(t, app)
	runReadyMetricsSession(t, app)
	require.Equal(t, int32(2), callbackCount.Load())
	require.Empty(t, app.automationCatalogSnapshot())

	evaluations, err := app.metricsRuntime.recorder.Reader().Evaluations(context.Background(), metrics.EvaluationFilter{})
	require.NoError(t, err)
	require.Empty(t, evaluations)
	app.Cleanup()
}
