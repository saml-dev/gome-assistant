package gomeassistant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
	"saml.dev/gome-assistant/internal/services"
)

func TestMetricsDisabledKeepsOriginalStateObject(t *testing.T) {
	app := newConfigurationTestApp(t)
	evaluation := app.newEvaluation("untraced", dailyTrigger())

	state := evaluation.scopedState(app.state)
	require.Same(t, app.state, state)
}

func TestScopedStateEqualsRecordsOneGetOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, "/api/states/sensor.test", request.URL.Path)
		_, _ = writer.Write([]byte(`{"entity_id":"sensor.test","state":"on","attributes":{}}`))
	}))
	t.Cleanup(server.Close)

	app, err := NewApp(context.Background(), NewAppRequest{URL: server.URL, HAAuthToken: "token"})
	require.NoError(t, err)
	recorder := openScopedTestRecorder(t)
	app.metricsOrigin.Store(&metricsSession{recorder: recorder, sessionID: "session"})

	evaluation := app.newEvaluation("daily", dailyTrigger())
	state := evaluation.scopedState(app.state)
	matches, err := state.Equals("sensor.test", "on")
	require.NoError(t, err)
	require.True(t, matches)

	events := waitForRecordedEvents(t, recorder, 3)
	require.Equal(t, []metrics.EventType{
		metrics.EventAutomationTriggered,
		metrics.EventStateReadStarted,
		metrics.EventStateReadFinished,
	}, oldestFirst(events))
	require.Equal(t, events[1].OperationID, events[2].OperationID)
	require.NotEmpty(t, events[1].OperationID)

	var started struct {
		Method  string         `json:"method"`
		Request map[string]any `json:"request"`
	}
	require.NoError(t, json.Unmarshal(events[1].PayloadJSON, &started))
	require.Equal(t, "get", started.Method)
	require.Equal(t, "sensor.test", started.Request["entity_id"])
}

func TestScopedServiceRecordsSynchronousOperationsOnly(t *testing.T) {
	recorder := openScopedTestRecorder(t)
	app := newConfigurationTestApp(t)
	app.metricsOrigin.Store(&metricsSession{recorder: recorder, sessionID: "session"})
	evaluation := app.newEvaluation("service", dailyTrigger())
	delegate := &scopedTestServicesAPI{}
	scopedAPI := &scopedServicesAPI{delegate: delegate, evaluation: evaluation}
	service := newService(scopedAPI)

	result, err := service.HomeAssistant.TurnOn(context.Background(), services.Entities("light.test"), map[string]any{"brightness": 20})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"ok": true}, result)
	require.NoError(t, service.Event.Fire("custom.event", map[string]any{"answer": 42}))
	require.NoError(t, scopedAPI.CallAndForget(services.BaseServiceRequest{Domain: "light", Service: "turn_off"}))
	require.Equal(t, int32(1), delegate.callAndForget.Load())

	events := waitForRecordedEvents(t, recorder, 5)
	require.Equal(t, []metrics.EventType{
		metrics.EventAutomationTriggered,
		metrics.EventServiceCallStarted,
		metrics.EventServiceCallFinished,
		metrics.EventServiceCallStarted,
		metrics.EventServiceCallFinished,
	}, oldestFirst(events))
	require.Equal(t, events[1].OperationID, events[2].OperationID)
	require.Equal(t, events[3].OperationID, events[4].OperationID)
	require.NotEqual(t, events[1].OperationID, events[3].OperationID)
}

type scopedTestServicesAPI struct {
	callAndForget atomic.Int32
}

func (api *scopedTestServicesAPI) CallAndForget(services.BaseServiceRequest) error {
	api.callAndForget.Add(1)
	return nil
}

func (api *scopedTestServicesAPI) Call(_ context.Context, _ services.BaseServiceRequest, result any) error {
	*result.(*any) = map[string]any{"ok": true}
	return nil
}

func (api *scopedTestServicesAPI) FireEvent(string, map[string]any) error {
	return nil
}

func openScopedTestRecorder(t *testing.T) *metrics.Recorder {
	t.Helper()
	recorder, err := metrics.OpenRecorder(filepath.Join(t.TempDir(), "metrics.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, recorder.Close(nil)) })
	return recorder
}

func waitForRecordedEvents(t *testing.T, recorder *metrics.Recorder, count int) []metrics.StoredEvent {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for {
		events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{Limit: 100, OldestFirst: true})
		require.NoError(t, err)
		if len(events) >= count {
			return events
		}
		if time.Now().After(deadline) {
			t.Fatalf("recorded %d events; wanted at least %d", len(events), count)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func oldestFirst(events []metrics.StoredEvent) []metrics.EventType {
	types := make([]metrics.EventType, len(events))
	for index, event := range events {
		types[index] = event.Type
	}
	return types
}
