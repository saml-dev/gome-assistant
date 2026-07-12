package gomeassistant

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEntityListenerCopiesShareDelayState(t *testing.T) {
	app, _ := newTestApp()
	app.entityListeners = make(map[string][]*EntityListener)
	var calls atomic.Int32
	listener := NewEntityListener().EntityIDs("light.test").Call(func(*Service, State, EntityData) {
		calls.Add(1)
	}).ToState("on").Duration("30ms").Build()
	listenerCopy := listener
	listenerCopy.entityIDs = []string{"light.test", "light.other"}
	app.RegisterEntityListeners(listenerCopy)

	app.callEntityListeners([]byte(`{"event":{"data":{"entity_id":"light.test","old_state":{"state":"off"},"new_state":{"state":"on"}}}}`))
	app.callEntityListeners([]byte(`{"event":{"data":{"entity_id":"light.other","old_state":{"state":"off"},"new_state":{"state":"on"}}}}`))

	time.Sleep(60 * time.Millisecond)
	require.EqualValues(t, 1, calls.Load())
	app.Cleanup()
}
