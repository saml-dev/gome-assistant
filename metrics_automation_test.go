package gomeassistant

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
	"saml.dev/gome-assistant/internal/services"
	"saml.dev/gome-assistant/websocket"
)

func TestMetricsAutomationTypesRecordSuccessAndSkip(t *testing.T) {
	t.Run("daily", func(t *testing.T) {
		app, recorder := newDirectMetricsTestApp(t)
		called := make(chan struct{})
		success := NewDailySchedule().Call(func(*Service, State) { close(called) }).At("12:00").Build()
		success.automationKey = "daily-success"
		success.maybeRunCallback(app)
		requireCallback(t, called)
		requireAutomationEvents(t, recorder, "daily-success", metrics.EventAutomationTriggered, metrics.EventAutomationCallbackStarted, metrics.EventAutomationCallbackFinished)

		skipped := NewDailySchedule().Call(func(*Service, State) {}).At("12:00").ExceptionDates(time.Now()).Build()
		skipped.automationKey = "daily-skipped"
		skipped.maybeRunCallback(app)
		requireSkippedAutomation(t, recorder, "daily-skipped", "exception_date")
	})

	t.Run("interval", func(t *testing.T) {
		app, recorder := newDirectMetricsTestApp(t)
		called := make(chan struct{})
		success := NewInterval().Call(func(*Service, State) { close(called) }).Every("1h").Build()
		success.automationKey = "interval-success"
		success.maybeRunCallback(app)
		requireCallback(t, called)
		requireAutomationEvents(t, recorder, "interval-success", metrics.EventAutomationTriggered, metrics.EventAutomationCallbackStarted, metrics.EventAutomationCallbackFinished)

		skipped := NewInterval().Call(func(*Service, State) {}).Every("1h").ExceptionDates(time.Now()).Build()
		skipped.automationKey = "interval-skipped"
		skipped.maybeRunCallback(app)
		requireSkippedAutomation(t, recorder, "interval-skipped", "exception_date")
	})

	t.Run("entity", func(t *testing.T) {
		app, recorder := newDirectMetricsTestApp(t)
		called := make(chan struct{})
		success := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) { close(called) }).Build()
		success.automationKey = "entity-success"
		data := stateData{EntityID: "light.test", OldState: msgState{State: "off"}, NewState: msgState{State: "on"}}
		success.maybeCall(app, EntityData{TriggerEntityID: "light.test", FromState: "off", ToState: "on"}, data)
		requireCallback(t, called)
		requireAutomationEvents(t, recorder, "entity-success", metrics.EventAutomationTriggered, metrics.EventAutomationCallbackStarted, metrics.EventAutomationCallbackFinished)

		skipped := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).IgnoreAttributeChanges().Build()
		skipped.automationKey = "entity-skipped"
		unchanged := stateData{EntityID: "light.test", OldState: msgState{State: "on"}, NewState: msgState{State: "on"}}
		skipped.maybeCall(app, EntityData{TriggerEntityID: "light.test", FromState: "on", ToState: "on"}, unchanged)
		requireSkippedAutomation(t, recorder, "entity-skipped", "attribute_only_change")
	})

	t.Run("event", func(t *testing.T) {
		app, recorder := newDirectMetricsTestApp(t)
		called := make(chan struct{})
		success := NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) { close(called) }).Build()
		success.automationKey = "event-success"
		success.maybeCall(app, EventData{Type: "test", RawEventJSON: []byte(`{"type":"event"}`)})
		requireCallback(t, called)
		requireAutomationEvents(t, recorder, "event-success", metrics.EventAutomationTriggered, metrics.EventAutomationCallbackStarted, metrics.EventAutomationCallbackFinished)

		skipped := NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) {}).ExceptionDates(time.Now()).Build()
		skipped.automationKey = "event-skipped"
		skipped.maybeCall(app, EventData{Type: "test", RawEventJSON: []byte(`{"type":"event"}`)})
		requireSkippedAutomation(t, recorder, "event-skipped", "exception_date")
	})
}

func TestMalformedMessagesDoNotCreateEvaluations(t *testing.T) {
	app, recorder := newDirectMetricsTestApp(t)
	entity := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).Build()
	entity.automationKey = "entity"
	app.entityListeners["light.test"] = []*EntityListener{&entity}
	event := NewEventListener().EventTypes("test").Call(func(*Service, State, EventData) {}).Build()
	event.automationKey = "event"
	app.eventListeners["test"] = []*EventListener{&event}

	app.callEntityListeners([]byte(`{"event":`))
	app.callEntityListeners([]byte(`{"type":"result","event":{"data":{"entity_id":"light.test"}}}`))
	app.callEventListeners("test", websocketMessage("result"))
	time.Sleep(20 * time.Millisecond)

	events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{Limit: 10})
	require.NoError(t, err)
	require.Empty(t, events)
}

func newDirectMetricsTestApp(t *testing.T) (*App, *metrics.Recorder) {
	t.Helper()
	app := newConfigurationTestApp(t)
	app.beginWorkerAdmission()
	app.setSessionContext(context.Background())
	recorder := openScopedTestRecorder(t)
	app.metricsOrigin.Store(&metricsSession{recorder: recorder, sessionID: "session"})
	t.Cleanup(func() {
		app.endWorkerAdmission()
		app.workers.Wait()
	})
	return app, recorder
}

func requireCallback(t *testing.T, called <-chan struct{}) {
	t.Helper()
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("automation callback was not invoked")
	}
}

func requireAutomationEvents(t *testing.T, recorder *metrics.Recorder, key string, expected ...metrics.EventType) []metrics.StoredEvent {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for {
		events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{
			AutomationKey: key,
			Limit:         100,
			OldestFirst:   true,
		})
		require.NoError(t, err)
		if len(events) >= len(expected) {
			require.Equal(t, expected, oldestFirst(events))
			return events
		}
		if time.Now().After(deadline) {
			t.Fatalf("recorded %d events for %q; wanted %d", len(events), key, len(expected))
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func requireSkippedAutomation(t *testing.T, recorder *metrics.Recorder, key, reason string) {
	t.Helper()
	events := requireAutomationEvents(t, recorder, key, metrics.EventAutomationTriggered, metrics.EventAutomationSkipped)
	var payload metrics.AutomationSkippedPayload
	require.NoError(t, json.Unmarshal(events[1].PayloadJSON, &payload))
	require.Equal(t, reason, payload.Reason)
}

// websocketMessage keeps malformed-message setup focused on the public shape
// consumed by callEventListeners.
func websocketMessage(messageType string) websocket.Message {
	return websocket.Message{BaseMessage: websocket.BaseMessage{Type: messageType}}
}

func TestEvaluationScopedObjectsRemainAttributedAfterCallbackReturn(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	databasePath := filepath.Join(t.TempDir(), "metrics.db")

	releaseChild := make(chan struct{})
	callbackReturned := make(chan struct{})
	childDone := make(chan error, 1)
	app.RegisterEntityListeners(NewEntityListener().
		EntityIDs("light.callback").
		Call(func(service *Service, state State, _ EntityData) {
			go func() {
				<-releaseChild
				if _, err := state.Get("light.child"); err != nil {
					childDone <- err
					return
				}
				_, err := service.HomeAssistant.TurnOn(context.Background(), services.Entities("light.child"))
				childDone <- err
			}()
			close(callbackReturned)
		}).
		RunOnStartup().
		Key("child-attribution").
		Build())
	require.NoError(t, app.EnableMetricsCollection(databasePath))

	runCtx, cancel := context.WithCancel(context.Background())
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(runCtx) }()
	t.Cleanup(func() {
		cancel()
		app.Cleanup()
	})
	requireCallback(t, callbackReturned)

	var evaluationID string
	deadline := time.Now().Add(time.Second)
	for evaluationID == "" {
		app.metricsMu.Lock()
		runtime := app.metricsRuntime
		app.metricsMu.Unlock()
		if runtime != nil {
			evaluations, err := runtime.recorder.Reader().Evaluations(context.Background(), metrics.EvaluationFilter{AutomationKey: "child-attribution"})
			require.NoError(t, err)
			if len(evaluations) == 1 {
				trace, err := runtime.recorder.Reader().EvaluationTrace(context.Background(), evaluations[0].EvaluationID)
				require.NoError(t, err)
				for _, event := range trace {
					if event.Type == metrics.EventAutomationCallbackFinished {
						evaluationID = evaluations[0].EvaluationID
						break
					}
				}
			}
		}
		if time.Now().After(deadline) {
			t.Fatal("callback-finished event was not committed")
		}
		time.Sleep(5 * time.Millisecond)
	}

	close(releaseChild)
	require.NoError(t, <-childDone)

	app.metricsMu.Lock()
	runtime := app.metricsRuntime
	app.metricsMu.Unlock()
	require.NotNil(t, runtime)
	deadline = time.Now().Add(time.Second)
	for {
		trace, err := runtime.recorder.Reader().EvaluationTrace(context.Background(), evaluationID)
		require.NoError(t, err)
		if len(trace) >= 9 {
			require.Equal(t, []metrics.EventType{
				metrics.EventAutomationTriggered,
				metrics.EventStateReadStarted,
				metrics.EventStateReadFinished,
				metrics.EventAutomationCallbackStarted,
				metrics.EventAutomationCallbackFinished,
				metrics.EventStateReadStarted,
				metrics.EventStateReadFinished,
				metrics.EventServiceCallStarted,
				metrics.EventServiceCallFinished,
			}, oldestFirst(trace))
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("recorded %d child-attribution events; wanted 9", len(trace))
		}
		time.Sleep(5 * time.Millisecond)
	}

	cancel()
	require.ErrorIs(t, <-runDone, context.Canceled)
}

func TestUnscopedStateAndServiceCallsDoNotCreateOperationEvents(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	runCtx, cancel := context.WithCancel(context.Background())
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(runCtx) }()
	deadline := time.Now().Add(time.Second)
	for !app.snapshotRuntimeStatus().SessionReady {
		if time.Now().After(deadline) {
			t.Fatal("metrics session did not become ready")
		}
		time.Sleep(5 * time.Millisecond)
	}

	_, err := app.GetState().Get("light.unscoped")
	require.NoError(t, err)
	_, err = app.GetService().HomeAssistant.TurnOn(context.Background(), services.Entities("light.unscoped"))
	require.NoError(t, err)
	time.Sleep(20 * time.Millisecond)

	events, err := app.metricsRuntime.recorder.Reader().Events(context.Background(), metrics.EventFilter{
		EventTypes: []metrics.EventType{
			metrics.EventStateReadStarted,
			metrics.EventStateReadFinished,
			metrics.EventServiceCallStarted,
			metrics.EventServiceCallFinished,
		},
		Limit: 100,
	})
	require.NoError(t, err)
	require.Empty(t, events)

	cancel()
	require.ErrorIs(t, <-runDone, context.Canceled)
	app.Cleanup()
}
