package gomeassistant

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
	"saml.dev/gome-assistant/websocket"
)

func TestEventListenerTriggerWhenIsOptionalAndRejectsNil(t *testing.T) {
	withoutPredicate := NewEventListener().
		EventTypes("test").
		Call(func(*Service, State, EventData) {}).
		Build()
	require.Nil(t, withoutPredicate.triggerWhen)

	withPredicateBuilder := NewEventListener().
		EventTypes("test").
		TriggerWhen(func(EventData) bool { return true })
	_, canSetAgain := reflect.TypeOf(withPredicateBuilder).MethodByName("TriggerWhen")
	require.False(t, canSetAgain)

	withPredicate := withPredicateBuilder.
		Call(func(*Service, State, EventData) {}).
		Build()
	require.NotNil(t, withPredicate.triggerWhen)

	require.PanicsWithValue(t, ErrInvalidArgs, func() {
		NewEventListener().EventTypes("test").TriggerWhen(nil)
	})
}

func TestEventListenerTriggerWhenFalseSuppressesEvaluationAndCallback(t *testing.T) {
	app, recorder := newDirectMetricsTestApp(t)
	var called atomic.Bool
	eventData := EventData{Type: "test", RawEventJSON: []byte(`{"accepted":false}`)}

	listener := NewEventListener().
		EventTypes("test").
		TriggerWhen(func(data EventData) bool {
			require.Equal(t, eventData, data)
			return false
		}).
		Call(func(*Service, State, EventData) { called.Store(true) }).
		ExceptionDates(time.Now()).
		Build()
	listener.automationKey = "trigger-when-false"
	lastRan := listener.lastRun()

	listener.maybeCall(app, eventData)
	app.workers.Wait()

	require.False(t, called.Load())
	require.Equal(t, lastRan, listener.lastRun())
	events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{
		AutomationKey: listener.automationKey,
		Limit:         10,
	})
	require.NoError(t, err)
	require.Empty(t, events)
}

func TestEventListenerTriggerWhenTruePreservesMetricsBehavior(t *testing.T) {
	app, recorder := newDirectMetricsTestApp(t)
	called := make(chan EventData, 1)
	eventData := EventData{Type: "test", RawEventJSON: []byte(`{"accepted":true}`)}

	listener := NewEventListener().
		EventTypes("test").
		TriggerWhen(func(data EventData) bool { return data.Type == "test" }).
		Call(func(_ *Service, _ State, data EventData) { called <- data }).
		Build()
	listener.automationKey = "trigger-when-true"

	listener.maybeCall(app, eventData)
	require.Equal(t, eventData, requireEventListenerData(t, called))
	requireAutomationEvents(t, recorder, listener.automationKey,
		metrics.EventAutomationTriggered,
		metrics.EventAutomationCallbackStarted,
		metrics.EventAutomationCallbackFinished,
	)
}

func TestEventListenerTriggerWhenAcceptsSubscribedEventWithoutMetrics(t *testing.T) {
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)
	t.Cleanup(func() { app.workers.Wait() })

	called := make(chan EventData, 1)
	raw := []byte(`{"event":{"data":{"value":42}}}`)
	listener := NewEventListener().
		EventTypes("accepted_event").
		TriggerWhen(func(data EventData) bool {
			return data.Type == "accepted_event" && string(data.RawEventJSON) == string(raw)
		}).
		Call(func(_ *Service, _ State, data EventData) { called <- data }).
		Build()
	app.RegisterEventListeners(listener)

	app.callEventListeners("accepted_event", websocket.Message{
		BaseMessage: websocket.BaseMessage{Type: "event"},
		Raw:         raw,
	})

	require.Equal(t, EventData{Type: "accepted_event", RawEventJSON: raw}, requireEventListenerData(t, called))
}

func TestEventListenerCallbackAdmissionFailureDoesNotUpdateLastRan(t *testing.T) {
	app, _ := newDirectMetricsTestApp(t)
	called := make(chan struct{}, 1)
	var attempts atomic.Int32

	listener := NewEventListener().
		EventTypes("test").
		TriggerWhen(func(EventData) bool {
			if attempts.Add(1) == 1 {
				app.endWorkerAdmission()
			}
			return true
		}).
		Call(func(*Service, State, EventData) { called <- struct{}{} }).
		Throttle("1h").
		Build()
	listener.automationKey = "trigger-when-admission"
	lastRan := listener.lastRun()

	listener.maybeCall(app, EventData{Type: "test"})
	require.Equal(t, lastRan, listener.lastRun())
	select {
	case <-called:
		t.Fatal("callback ran after admission was closed")
	default:
	}

	app.beginWorkerAdmission()
	listener.maybeCall(app, EventData{Type: "test"})
	requireCallback(t, called)
	require.NotEqual(t, lastRan, listener.lastRun())
}

func requireEventListenerData(t *testing.T, called <-chan EventData) EventData {
	t.Helper()
	select {
	case data := <-called:
		return data
	case <-time.After(time.Second):
		t.Fatal("event listener callback was not invoked")
		return EventData{}
	}
}
