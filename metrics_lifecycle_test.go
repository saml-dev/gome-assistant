package gomeassistant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"saml.dev/gome-assistant/internal/metrics"
)

func TestMetricsRuntimePersistsAcrossSequentialSessionsAndClosesOnCleanup(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))

	firstSession := runReadyMetricsSession(t, app)
	runtime := app.metricsRuntime
	require.NotNil(t, runtime)
	require.True(t, runtime.promoted.Load())

	secondSession := runReadyMetricsSession(t, app)
	require.Same(t, runtime, app.metricsRuntime)
	require.NotEqual(t, firstSession, secondSession)

	events := waitForRecordedEvents(t, runtime.recorder, 6)
	starts, readies, finishes := 0, 0, 0
	sessions := make(map[string]struct{})
	for _, event := range events {
		switch event.Type {
		case metrics.EventSessionStarting:
			starts++
			sessions[event.SessionID] = struct{}{}
		case metrics.EventSessionReady:
			readies++
		case metrics.EventSessionFinished:
			finishes++
		}
	}
	require.Equal(t, 2, starts)
	require.Equal(t, 2, readies)
	require.Equal(t, 2, finishes)
	require.Equal(t, map[string]struct{}{firstSession: {}, secondSession: {}}, sessions)

	app.Cleanup()
	require.True(t, runtime.recorder.Health().Closed)
	require.True(t, app.snapshotRuntimeStatus().PermanentlyClosed)
}

func TestFailedInitialMetricsSessionIsProvisionalAndRecordsStartupFailure(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(server.Close)

	app := fakeApp(t, server.URL)
	databasePath := filepath.Join(t.TempDir(), "metrics.db")
	require.NoError(t, app.EnableMetricsCollection(databasePath))
	require.Error(t, app.Start(context.Background()))
	require.Nil(t, app.metricsRuntime, "an initial failed session must not retain the runtime")

	recorder, err := metrics.OpenRecorder(databasePath)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, recorder.Close(nil)) })
	events := waitForRecordedEvents(t, recorder, 2)
	require.Equal(t, []metrics.EventType{
		metrics.EventSessionStarting,
		metrics.EventSessionFinished,
	}, oldestFirst(events))

	var finished struct {
		Reason string         `json:"reason"`
		Error  *metrics.Error `json:"error"`
	}
	require.NoError(t, json.Unmarshal(events[1].PayloadJSON, &finished))
	require.Equal(t, "startup_error", finished.Reason)
	require.NotNil(t, finished.Error)
}

func TestCleanupBeforeLaterSessionMetricsAdmissionClosesPromotedRuntime(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)
	require.NoError(t, app.EnableMetricsCollection(filepath.Join(t.TempDir(), "metrics.db")))
	_ = runReadyMetricsSession(t, app)
	runtime := app.metricsRuntime
	require.NotNil(t, runtime)
	require.True(t, runtime.promoted.Load())

	// Hold Start at worker admission after it has published its active run,
	// but before its first pre-metrics context check. Cleanup therefore makes
	// this invocation return with a nil metricsSession.
	app.workerAdmissionMu.Lock()
	done := make(chan error, 1)
	go func() { done <- app.Start(context.Background()) }()
	deadline := time.Now().Add(time.Second)
	for app.run.Load() == nil {
		if time.Now().After(deadline) {
			app.workerAdmissionMu.Unlock()
			t.Fatal("later Start did not publish its cancellation handle")
		}
		time.Sleep(time.Millisecond)
	}
	app.Cleanup()
	app.workerAdmissionMu.Unlock()

	require.ErrorIs(t, <-done, context.Canceled)
	require.True(t, runtime.recorder.Health().Closed)
	require.True(t, app.snapshotRuntimeStatus().PermanentlyClosed)
}

func runReadyMetricsSession(t *testing.T, app *App) string {
	t.Helper()
	runCtx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- app.Start(runCtx) }()

	deadline := time.Now().Add(time.Second)
	for !app.snapshotRuntimeStatus().SessionReady {
		if time.Now().After(deadline) {
			t.Fatal("metrics session did not become ready")
		}
		time.Sleep(5 * time.Millisecond)
	}
	sessionID := app.snapshotRuntimeStatus().SessionID
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	return sessionID
}
