package gomeassistant

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/carbon"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/websocket"
)

type EventListener struct {
	eventTypes   []string
	triggerWhen  func(EventData) bool
	callback     EventListenerCallback
	betweenStart string
	betweenEnd   string
	throttle     time.Duration
	lastRan      carbon.Carbon
	lastRanMu    *sync.Mutex

	exceptionDates  []time.Time
	exceptionRanges []timeRange

	enabledEntities  []internal.EnabledDisabledInfo
	disabledEntities []internal.EnabledDisabledInfo

	key             string
	automationKey   string
	source          string
	metricsDisabled bool
}

type EventListenerCallback func(*Service, State, EventData)

type EventData struct {
	Type         string
	RawEventJSON []byte
}

/* Methods */

func NewEventListener() eventListenerBuilder1 {
	return eventListenerBuilder1{EventListener{
		lastRan:   carbon.Now().StartOfCentury(),
		lastRanMu: &sync.Mutex{},
	}}
}

type eventListenerBuilder1 struct {
	eventListener EventListener
}

func (b eventListenerBuilder1) EventTypes(ets ...string) eventListenerBuilder2 {
	b.eventListener.eventTypes = ets
	return eventListenerBuilder2(b)
}

type eventListenerBuilder2 struct {
	eventListener EventListener
}

// TriggerWhen refines whether an event type match triggers this listener. A
// false result excludes the event before evaluation, so it produces no
// metrics, condition checks, callback, or throttle update. It may be set once,
// before Call.
func (b eventListenerBuilder2) TriggerWhen(predicate func(EventData) bool) eventListenerBuilderTriggerWhen {
	if predicate == nil {
		panic(ErrInvalidArgs)
	}
	b.eventListener.triggerWhen = predicate
	return eventListenerBuilderTriggerWhen(b)
}

func (b eventListenerBuilder2) Call(callback EventListenerCallback) eventListenerBuilder3 {
	b.eventListener.callback = callback
	return eventListenerBuilder3(b)
}

type eventListenerBuilderTriggerWhen struct {
	eventListener EventListener
}

func (b eventListenerBuilderTriggerWhen) Call(callback EventListenerCallback) eventListenerBuilder3 {
	b.eventListener.callback = callback
	return eventListenerBuilder3(b)
}

type eventListenerBuilder3 struct {
	eventListener EventListener
}

func (b eventListenerBuilder3) OnlyBetween(start string, end string) eventListenerBuilder3 {
	b.eventListener.betweenStart = start
	b.eventListener.betweenEnd = end
	return b
}

func (b eventListenerBuilder3) OnlyAfter(start string) eventListenerBuilder3 {
	b.eventListener.betweenStart = start
	return b
}

func (b eventListenerBuilder3) OnlyBefore(end string) eventListenerBuilder3 {
	b.eventListener.betweenEnd = end
	return b
}

func (b eventListenerBuilder3) Throttle(s DurationString) eventListenerBuilder3 {
	d := internal.ParseDuration(string(s))
	b.eventListener.throttle = d
	return b
}

func (b eventListenerBuilder3) ExceptionDates(t time.Time, tl ...time.Time) eventListenerBuilder3 {
	b.eventListener.exceptionDates = append(tl, t)
	return b
}

func (b eventListenerBuilder3) ExceptionRange(start, end time.Time) eventListenerBuilder3 {
	b.eventListener.exceptionRanges = append(b.eventListener.exceptionRanges, timeRange{start, end})
	return b
}

/*
Enable this listener only when the current state of {entityID} matches {state}.
If there is a network error while retrieving state, the listener runs if {runOnNetworkError} is true.
*/
func (b eventListenerBuilder3) EnabledWhen(entityID, state string, runOnNetworkError bool) eventListenerBuilder3 {
	if entityID == "" {
		panic(fmt.Sprintf("entityID is empty in eventListener EnabledWhen entityID='%s' state='%s' runOnNetworkError='%t'", entityID, state, runOnNetworkError))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.eventListener.enabledEntities = append(b.eventListener.enabledEntities, i)
	return b
}

/*
Disable this listener when the current state of {entityID} matches {state}.
If there is a network error while retrieving state, the listener runs if {runOnNetworkError} is true.
*/
func (b eventListenerBuilder3) DisabledWhen(entityID, state string, runOnNetworkError bool) eventListenerBuilder3 {
	if entityID == "" {
		panic(fmt.Sprintf("entityID is empty in eventListener EnabledWhen entityID='%s' state='%s' runOnNetworkError='%t'", entityID, state, runOnNetworkError))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.eventListener.disabledEntities = append(b.eventListener.disabledEntities, i)
	return b
}

// Key sets the stable logical identity used by metrics and the web UI.
func (b eventListenerBuilder3) Key(key string) eventListenerBuilder3 {
	key = strings.TrimSpace(key)
	if key == "" {
		panic(ErrInvalidArgs)
	}
	b.eventListener.key = key
	return b
}

// DisableMetrics excludes this listener from metrics collection and the web
// UI catalog without disabling its execution.
func (b eventListenerBuilder3) DisableMetrics() eventListenerBuilder3 {
	b.eventListener.metricsDisabled = true
	return b
}

func (b eventListenerBuilder3) Build() EventListener {
	return b.eventListener
}

func (app *App) callbackAdmissionOpen() bool {
	return !app.sessionStopping()
}

func (l *EventListener) maybeCall(app *App, eventData EventData) {
	if l.triggerWhen != nil && !l.triggerWhen(eventData) {
		return
	}

	evaluation := app.newEvaluation(l.automationKey, eventTrigger(eventData))
	if !app.callbackAdmissionOpen() {
		evaluation.skip("session_stopping", nil, nil)
		return
	}

	// Check conditions
	if c := checkWithinTimeRange(l.betweenStart, l.betweenEnd); c.fail {
		evaluation.skip("outside_time_window", nil, nil)
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

	if app.sessionStopping() {
		evaluation.skip("session_stopping", nil, nil)
		return
	}
	if !evaluation.runEventCallback(app, l.callback, eventData) {
		evaluation.skip("session_stopping", nil, nil)
		return
	}
	l.markRan()
}

func (l *EventListener) lastRun() carbon.Carbon {
	mutex := l.lastRanMutex()
	mutex.Lock()
	defer mutex.Unlock()
	return l.lastRan
}

func (l *EventListener) lastRanMutex() *sync.Mutex {
	if l.lastRanMu != nil {
		return l.lastRanMu
	}
	return &fallbackLastRanMutex
}

func (l *EventListener) markRan() {
	mutex := l.lastRanMutex()
	mutex.Lock()
	l.lastRan = carbon.Now()
	mutex.Unlock()
}

/* Functions */
func (app *App) callEventListeners(eventType string, msg websocket.Message) {
	if msg.Type != "event" {
		return
	}
	listeners, ok := app.eventListeners[eventType]
	if !ok {
		// no listeners registered for this event type
		return
	}

	eventData := EventData{
		Type:         eventType,
		RawEventJSON: msg.Raw,
	}

	for _, l := range listeners {
		l.maybeCall(app, eventData)
	}
}
