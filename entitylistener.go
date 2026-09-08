package gomeassistant

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/carbon"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/websocket"
)

type EntityListener struct {
	entityIDs []string
	callback  EntityListenerCallback
	fromState string
	toState   string
	throttle  time.Duration
	lastRan   carbon.Carbon
	lastRanMu *sync.Mutex

	betweenStart string
	betweenEnd   string

	delay      time.Duration
	delayState *entityDelayState

	exceptionDates  []time.Time
	exceptionRanges []timeRange

	runOnStartup bool

	ignoreAttributeChanges bool

	enabledEntities  []internal.EnabledDisabledInfo
	disabledEntities []internal.EnabledDisabledInfo

	key             string
	automationKey   string
	source          string
	metricsDisabled bool
}

type entityDelayWait struct {
	timer      *time.Timer
	cancel     context.CancelFunc
	evaluation *evaluationTrace
}

type entityDelayState struct {
	sync.Mutex
	delayWait *entityDelayWait
}

type EntityListenerCallback func(*Service, State, EntityData)

var fallbackLastRanMutex sync.Mutex

type EntityData struct {
	TriggerEntityID string
	FromState       string
	FromAttributes  map[string]any
	ToState         string
	ToAttributes    map[string]any
	LastChanged     time.Time
}

type stateChangedMessage struct {
	websocket.BaseMessage
	Event struct {
		Data      stateData `json:"data"`
		EventType string    `json:"event_type"`
		Origin    string    `json:"origin"`
	} `json:"event"`
}

type stateData struct {
	EntityID string   `json:"entity_id"`
	NewState msgState `json:"new_state"`
	OldState msgState `json:"old_state"`
}

type msgState struct {
	EntityID    string         `json:"entity_id"`
	LastChanged time.Time      `json:"last_changed"`
	State       string         `json:"state"`
	Attributes  map[string]any `json:"attributes"`
}

/* Methods */

func NewEntityListener() elBuilder1 {
	return elBuilder1{EntityListener{
		lastRan:    carbon.Now().StartOfCentury(),
		lastRanMu:  &sync.Mutex{},
		delayState: &entityDelayState{},
	}}
}

type elBuilder1 struct {
	entityListener EntityListener
}

func (b elBuilder1) EntityIDs(EntityIDs ...string) elBuilder2 {
	if len(EntityIDs) == 0 {
		panic("must pass at least one entityID to EntityIDs()")
	} else {
		b.entityListener.entityIDs = EntityIDs
	}
	return elBuilder2(b)
}

type elBuilder2 struct {
	entityListener EntityListener
}

func (b elBuilder2) Call(callback EntityListenerCallback) elBuilder3 {
	b.entityListener.callback = callback
	return elBuilder3(b)
}

type elBuilder3 struct {
	entityListener EntityListener
}

func (b elBuilder3) OnlyBetween(start string, end string) elBuilder3 {
	b.entityListener.betweenStart = start
	b.entityListener.betweenEnd = end
	return b
}

func (b elBuilder3) OnlyAfter(start string) elBuilder3 {
	b.entityListener.betweenStart = start
	return b
}

func (b elBuilder3) OnlyBefore(end string) elBuilder3 {
	b.entityListener.betweenEnd = end
	return b
}

func (b elBuilder3) FromState(s string) elBuilder3 {
	b.entityListener.fromState = s
	return b
}

func (b elBuilder3) ToState(s string) elBuilder3 {
	b.entityListener.toState = s
	return b
}

func (b elBuilder3) Duration(s DurationString) elBuilder3 {
	d := internal.ParseDuration(string(s))
	b.entityListener.delay = d
	return b
}

func (b elBuilder3) Throttle(s DurationString) elBuilder3 {
	d := internal.ParseDuration(string(s))
	b.entityListener.throttle = d
	return b
}

func (b elBuilder3) ExceptionDates(t time.Time, tl ...time.Time) elBuilder3 {
	b.entityListener.exceptionDates = append(tl, t)
	return b
}

func (b elBuilder3) ExceptionRange(start, end time.Time) elBuilder3 {
	b.entityListener.exceptionRanges = append(b.entityListener.exceptionRanges, timeRange{start, end})
	return b
}

// RunOnStartup invokes the listener once after each successful Start session
// starts. A listener registered for multiple entities still runs only once per
// session.
func (b elBuilder3) RunOnStartup() elBuilder3 {
	b.entityListener.runOnStartup = true
	return b
}

// IgnoreAttributeChanges causes this listener to skip firing when the entity's
// state string is unchanged (i.e. only attributes changed). By default,
// attribute-only changes will trigger the listener.
func (b elBuilder3) IgnoreAttributeChanges() elBuilder3 {
	b.entityListener.ignoreAttributeChanges = true
	return b
}

/*
Enable this listener only when the current state of {entityID} matches {state}.
If there is a network error while retrieving state, the listener runs if {runOnNetworkError} is true.
*/
func (b elBuilder3) EnabledWhen(entityID, state string, runOnNetworkError bool) elBuilder3 {
	if entityID == "" {
		panic(fmt.Sprintf("entityID is empty in EnabledWhen entityID='%s' state='%s'", entityID, state))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.entityListener.enabledEntities = append(b.entityListener.enabledEntities, i)
	return b
}

/*
Disable this listener when the current state of {entityID} matches {state}.
If there is a network error while retrieving state, the listener runs if {runOnNetworkError} is true.
*/
func (b elBuilder3) DisabledWhen(entityID, state string, runOnNetworkError bool) elBuilder3 {
	if entityID == "" {
		panic(fmt.Sprintf("entityID is empty in EnabledWhen entityID='%s' state='%s'", entityID, state))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.entityListener.disabledEntities = append(b.entityListener.disabledEntities, i)
	return b
}

// Key sets the stable logical identity used by metrics and the web UI.
func (b elBuilder3) Key(key string) elBuilder3 {
	key = strings.TrimSpace(key)
	if key == "" {
		panic(ErrInvalidArgs)
	}
	b.entityListener.key = key
	return b
}

// DisableMetrics excludes this listener from metrics collection and the web
// UI catalog without disabling its execution.
func (b elBuilder3) DisableMetrics() elBuilder3 {
	b.entityListener.metricsDisabled = true
	return b
}

func (b elBuilder3) Build() EntityListener {
	return b.entityListener
}

func (l *EntityListener) maybeCall(app *App, entityData EntityData, data stateData) {
	evaluation := app.newEvaluation(l.automationKey, entityTrigger("entity", entityData))
	if !app.callbackAdmissionOpen() {
		evaluation.skip("session_stopping", nil, nil)
		return
	}

	// Check conditions
	if l.ignoreAttributeChanges && data.NewState.State == data.OldState.State {
		evaluation.skip("attribute_only_change", nil, nil)
		return
	}
	if c := checkWithinTimeRange(l.betweenStart, l.betweenEnd); c.fail {
		evaluation.skip("outside_time_window", nil, nil)
		return
	}
	if c := checkStatesMatch(l.fromState, data.OldState.State); c.fail {
		evaluation.skip("from_state_mismatch", nil, nil)
		return
	}
	if c := checkStatesMatch(l.toState, data.NewState.State); c.fail {
		l.cancelDelay("target_state_changed", evaluation)
		evaluation.skip("to_state_mismatch", nil, nil)
		return
	}
	if c := checkThrottle(l.throttle, l.lastRun()); c.fail {
		evaluation.skip("throttled", nil, nil)
		return
	}
	if c := checkExceptionDates(l.exceptionDates); c.fail {
		evaluation.skip("exception_date", nil, nil)
		return
	}
	if c := checkExceptionRanges(l.exceptionRanges); c.fail {
		evaluation.skip("exception_range", nil, nil)
		return
	}
	if c := checkEnabledEntity(evaluation.scopedState(app.state), l.enabledEntities); c.fail {
		evaluation.skipCondition(c, "enabled_condition_mismatch", "enabled_condition_query_error")
		return
	}
	if c := checkDisabledEntity(evaluation.scopedState(app.state), l.disabledEntities); c.fail {
		evaluation.skipCondition(c, "disabled_condition_match", "disabled_condition_query_error")
		return
	}

	if l.delay != 0 {
		l.scheduleDelay(app, entityData, evaluation)
		return
	}

	// run now if no delay set
	if app.sessionStopping() {
		evaluation.skip("session_stopping", nil, nil)
		return
	}
	if !evaluation.runEntityCallback(app, l.callback, entityData) {
		evaluation.skip("session_stopping", nil, nil)
	}
	l.markRan()
}

func (l *EntityListener) lastRun() carbon.Carbon {
	mutex := l.lastRanMutex()
	mutex.Lock()
	defer mutex.Unlock()
	return l.lastRan
}

func (l *EntityListener) lastRanMutex() *sync.Mutex {
	if l.lastRanMu != nil {
		return l.lastRanMu
	}
	return &fallbackLastRanMutex
}

func (l *EntityListener) markRan() {
	mutex := l.lastRanMutex()
	mutex.Lock()
	l.lastRan = carbon.Now()
	mutex.Unlock()
}

func (l *EntityListener) cancelDelay(reason string, cause *evaluationTrace) {
	l.delayState.Lock()
	wait := l.delayState.delayWait
	l.delayState.delayWait = nil
	l.delayState.Unlock()
	if wait != nil {
		wait.evaluation.cancelWaiting(reason, cause)
		wait.cancel()
	}
}

func (l *EntityListener) scheduleDelay(app *App, entityData EntityData, evaluation *evaluationTrace) {
	delayCtx, cancel := context.WithCancel(app.sessionContext())
	wait := &entityDelayWait{timer: time.NewTimer(l.delay), cancel: cancel, evaluation: evaluation}

	l.delayState.Lock()
	// Mark the trace waiting before its wait is visible to another trigger.
	// A replacement then always observes a cancelable active evaluation.
	evaluation.beginWaiting(l.delay, l.toState)
	previous := l.delayState.delayWait
	l.delayState.delayWait = wait
	l.delayState.Unlock()
	if previous != nil {
		previous.evaluation.cancelWaiting("replaced_by_new_trigger", evaluation)
		previous.cancel()
	}

	if !app.goTracked(func() {
		defer func() {
			wait.cancel()
			stopAndDrainTimer(wait.timer)
			l.delayState.Lock()
			if l.delayState.delayWait == wait {
				l.delayState.delayWait = nil
			}
			l.delayState.Unlock()
		}()

		select {
		case <-wait.timer.C:
		case <-delayCtx.Done():
			if app.sessionStopping() {
				evaluation.cancelWaiting("session_stopping", nil)
			} else {
				evaluation.cancelWaiting("stale_wait", nil)
			}
			return
		}
		if app.sessionStopping() || delayCtx.Err() != nil {
			evaluation.cancelWaiting("session_stopping", nil)
			return
		}

		l.delayState.Lock()
		if l.delayState.delayWait != wait {
			l.delayState.Unlock()
			evaluation.cancelWaiting("stale_wait", nil)
			return
		}
		l.delayState.delayWait = nil
		l.delayState.Unlock()
		if app.sessionStopping() || delayCtx.Err() != nil {
			evaluation.cancelWaiting("session_stopping", nil)
			return
		}
		evaluation.finishWaiting()
		if !evaluation.runEntityCallback(app, l.callback, entityData) {
			evaluation.skip("session_stopping", nil, nil)
		}
		l.markRan()
	}) {
		l.delayState.Lock()
		if l.delayState.delayWait == wait {
			l.delayState.delayWait = nil
		}
		l.delayState.Unlock()
		wait.cancel()
		stopAndDrainTimer(wait.timer)
		evaluation.cancelWaiting("session_stopping", nil)
	}
}

func stopAndDrainTimer(timer *time.Timer) {
	if timer.Stop() {
		return
	}
	select {
	case <-timer.C:
	default:
	}
}

/* Functions */
func (app *App) callEntityListeners(msgBytes []byte) {
	msg := stateChangedMessage{}
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return
	}
	if msg.Type != "" && msg.Type != "event" {
		return
	}
	data := msg.Event.Data
	eid := data.EntityID
	listeners, ok := app.entityListeners[eid]
	if !ok {
		// no listeners registered for this id
		return
	}

	entityData := EntityData{
		TriggerEntityID: eid,
		FromState:       data.OldState.State,
		FromAttributes:  data.OldState.Attributes,
		ToState:         data.NewState.State,
		ToAttributes:    data.NewState.Attributes,
		LastChanged:     data.OldState.LastChanged,
	}

	for _, l := range listeners {
		l.maybeCall(app, entityData, data)
	}
}
