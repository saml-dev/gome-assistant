package gomeassistant

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestApp() (*App, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	app := &App{ctx: ctx}
	app.run.Store(&runCancellation{cancel: cancel})
	return app, cancel
}

// Worker draining is the lifecycle primitive covered here; end-to-end
// Cleanup-with-active-Start coverage is exercised by the lifecycle acceptance
// tests.
func TestEntityListenerCallbackIsDrainedByLifecycleWorkerDrain(t *testing.T) {
	app, cancel := newTestApp()
	started := make(chan struct{})
	release := make(chan struct{})
	listener := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {
		close(started)
		<-release
	}).Build()

	listener.maybeCall(app, EntityData{}, stateData{
		OldState: msgState{State: "off"},
		NewState: msgState{State: "on"},
	})
	<-started

	drained := make(chan struct{})
	go func() {
		cancel()
		app.workers.Wait()
		close(drained)
	}()

	select {
	case <-drained:
		t.Fatal("lifecycle cancellation returned while callback was still running")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	select {
	case <-drained:
	case <-time.After(time.Second):
		t.Fatal("lifecycle cancellation did not wait for callback")
	}
}

func TestEntityListenerDelayCancellationPreventsReplacementCallback(t *testing.T) {
	app, _ := newTestApp()
	var calls atomic.Int32
	listener := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {
		calls.Add(1)
	}).ToState("on").Duration("30ms").Build()

	listener.maybeCall(app, EntityData{}, stateData{
		OldState: msgState{State: "off"},
		NewState: msgState{State: "on"},
	})
	listener.maybeCall(app, EntityData{}, stateData{
		OldState: msgState{State: "on"},
		NewState: msgState{State: "off"},
	})

	time.Sleep(60 * time.Millisecond)
	require.Zero(t, calls.Load())
	app.Cleanup()
}

func TestEntityListenerDelayCancellationPreventsCallbackAfterCleanup(t *testing.T) {
	app, _ := newTestApp()
	var calls atomic.Int32
	listener := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {
		calls.Add(1)
	}).Duration("30ms").Build()

	listener.maybeCall(app, EntityData{}, stateData{
		OldState: msgState{State: "off"},
		NewState: msgState{State: "on"},
	})
	app.Cleanup()
	time.Sleep(60 * time.Millisecond)
	require.Zero(t, calls.Load())
}
