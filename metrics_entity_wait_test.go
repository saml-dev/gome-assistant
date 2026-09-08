package gomeassistant

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
)

func TestConcurrentEntityDelayReplacementHasOneActiveWaitAndCancellation(t *testing.T) {
	app := newConfigurationTestApp(t)
	ctx, cancel := context.WithCancel(context.Background())
	app.setSessionContext(ctx)
	app.run.Store(&runCancellation{cancel: cancel})

	recorder := openScopedTestRecorder(t)
	t.Cleanup(func() {
		cancel()
		app.workers.Wait()
	})
	app.metricsOrigin.Store(&metricsSession{recorder: recorder, sessionID: "session"})
	listener := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {}).ToState("on").Duration("1h").Build()
	first := app.newEvaluation("entity", entityTrigger("entity", EntityData{TriggerEntityID: "light.test"}))
	second := app.newEvaluation("entity", entityTrigger("entity", EntityData{TriggerEntityID: "light.test"}))

	start := make(chan struct{})
	var workers sync.WaitGroup
	workers.Add(2)
	go func() {
		defer workers.Done()
		<-start
		listener.scheduleDelay(app, EntityData{TriggerEntityID: "light.test"}, first)
	}()
	go func() {
		defer workers.Done()
		<-start
		listener.scheduleDelay(app, EntityData{TriggerEntityID: "light.test"}, second)
	}()
	close(start)
	workers.Wait()

	events := waitForRecordedEvents(t, recorder, 5)
	waiting, canceled := 0, 0
	for _, event := range events {
		switch event.Type {
		case metrics.EventAutomationWaiting:
			waiting++
		case metrics.EventAutomationWaitCanceled:
			canceled++
			var payload metrics.AutomationWaitCanceledPayload
			require.NoError(t, json.Unmarshal(event.PayloadJSON, &payload))
			require.Equal(t, "replaced_by_new_trigger", payload.Reason)
			require.NotNil(t, payload.CauseEvaluationID)
			require.Contains(t, []string{first.evaluationID, second.evaluationID}, *payload.CauseEvaluationID)
		}
	}
	require.Equal(t, 2, waiting)
	require.Equal(t, 1, canceled)
	require.NotEqual(t, first.waiting.Load(), second.waiting.Load(), "exactly one replacement wait remains active")
}
