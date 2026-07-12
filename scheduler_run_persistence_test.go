package gomeassistant

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStartPersistsDailyScheduleNextRunTimeAcrossSessions(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)

	var calls atomic.Int32
	called := make(chan struct{}, 1)
	firstCtx, cancelFirstRun := context.WithCancel(context.Background())
	app.RegisterSchedules(NewDailySchedule().Call(func(*Service, State) {
		calls.Add(1)
		called <- struct{}{}
		cancelFirstRun()
	}).At("00:00").Build())

	schedule, ok := app.scheduledActions[0].(*DailySchedule)
	require.True(t, ok)
	schedule.nextRunTime = time.Now().Add(-time.Millisecond)

	startAndWaitForScheduledCall(t, app, firstCtx, called)
	require.Equal(t, int32(1), calls.Load())
	require.True(t, schedule.nextRunTime.After(time.Now()))

	startWithoutScheduledCall(t, app, called)
	require.Equal(t, int32(1), calls.Load())
}

func TestStartPersistsIntervalNextRunTimeAcrossSessions(t *testing.T) {
	server, _, _ := newFakeHA(t, http.StatusOK)
	app := fakeApp(t, server.URL)

	var calls atomic.Int32
	called := make(chan struct{}, 1)
	firstCtx, cancelFirstRun := context.WithCancel(context.Background())
	app.RegisterIntervals(NewInterval().Call(func(*Service, State) {
		calls.Add(1)
		called <- struct{}{}
		cancelFirstRun()
	}).Every("24h").Build())

	interval, ok := app.scheduledActions[0].(*Interval)
	require.True(t, ok)
	interval.nextRunTime = time.Now().Add(-time.Millisecond)

	startAndWaitForScheduledCall(t, app, firstCtx, called)
	require.Equal(t, int32(1), calls.Load())
	require.True(t, interval.nextRunTime.After(time.Now()))

	startWithoutScheduledCall(t, app, called)
	require.Equal(t, int32(1), calls.Load())
}

func TestScheduleAndIntervalRegistrationKeepInternalCopies(t *testing.T) {
	app, err := NewApp(context.Background(), NewAppRequest{
		URL:         "http://127.0.0.1:1",
		HAAuthToken: "token",
	})
	require.NoError(t, err)

	schedule := NewDailySchedule().Call(func(*Service, State) {}).At("00:00").Build()
	interval := NewInterval().Call(func(*Service, State) {}).Every("24h").Build()
	app.RegisterSchedules(schedule)
	app.RegisterIntervals(interval)

	registeredSchedule, ok := app.scheduledActions[0].(*DailySchedule)
	require.True(t, ok)
	registeredInterval, ok := app.scheduledActions[1].(*Interval)
	require.True(t, ok)
	require.NotSame(t, &schedule, registeredSchedule)
	require.NotSame(t, &interval, registeredInterval)
	require.True(t, schedule.nextRunTime.IsZero())
	require.True(t, interval.nextRunTime.IsZero())

	registeredScheduleNextRun := registeredSchedule.nextRunTime
	registeredIntervalNextRun := registeredInterval.nextRunTime
	schedule.nextRunTime = time.Unix(1, 0)
	interval.nextRunTime = time.Unix(1, 0)
	require.Equal(t, registeredScheduleNextRun, registeredSchedule.nextRunTime)
	require.Equal(t, registeredIntervalNextRun, registeredInterval.nextRunTime)
}

func startAndWaitForScheduledCall(t *testing.T, app *App, ctx context.Context, called <-chan struct{}) {
	t.Helper()
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(ctx) }()

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("scheduled callback was not invoked")
	}

	select {
	case err := <-runDone:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Start did not return after cancellation")
	}
}

func startWithoutScheduledCall(t *testing.T, app *App, called <-chan struct{}) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runDone := make(chan error, 1)
	go func() { runDone <- app.Start(ctx) }()

	select {
	case <-called:
		t.Fatal("scheduled callback replayed in the next Start session")
	case <-time.After(50 * time.Millisecond):
	}

	cancel()
	select {
	case err := <-runDone:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("Start did not return after cancellation")
	}
}
