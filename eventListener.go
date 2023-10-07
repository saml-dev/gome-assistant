package gomeassistant

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-module/carbon"
	"saml.dev/gome-assistant/internal"
	ws "saml.dev/gome-assistant/internal/websocket"
)

type EventListener struct {
	eventTypes   []string
	callback     EventListenerCallback
	betweenStart string
	betweenEnd   string
	throttle     time.Duration
	lastRan      carbon.Carbon

	exceptionDates  []time.Time
	exceptionRanges []timeRange

	enabledEntity            string
	enabledEntityState       string
	enabledEntityRunOnError  bool
	disabledEntity           string
	disabledEntityState      string
	disabledEntityRunOnError bool
}

type EventListenerCallback func(*Service, *State, EventData)

type EventData struct {
	Type         string
	RawEventJSON []byte
}

/* Methods */

func NewEventListener() eventListenerBuilder1 {
	return eventListenerBuilder1{EventListener{
		lastRan: carbon.Now().StartOfCentury(),
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

func (b eventListenerBuilder2) Call(callback EventListenerCallback) eventListenerBuilder3 {
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
Enable this listener only when the current state of {entityId} matches {state}.
If there is a network error while retrieving state, the listener runs if {runOnNetworkError} is true.
*/
func (b eventListenerBuilder3) EnabledWhen(entityId, state string, runOnNetworkError bool) eventListenerBuilder3 {
	if entityId == "" || state == "" {
		panic(fmt.Sprintf("Either entityId or state is empty in eventListener EnabledWhen entityId='%s' state='%s' runOnNetworkError='%t'", entityId, state, runOnNetworkError))
	}
	if b.eventListener.enabledEntity != "" {
		panic(fmt.Sprintf("You can't use EnabledWhen and DisabledWhen together. Error occurred while setting EnabledWhen entityId=%s state=%s runOnNetworkError=%t", entityId, state, runOnNetworkError))
	}
	b.eventListener.enabledEntity = entityId
	b.eventListener.enabledEntityState = state
	b.eventListener.enabledEntityRunOnError = runOnNetworkError
	return b
}

/*
Disable this listener when the current state of {entityId} matches {state}.
If there is a network error while retrieving state, the listener runs if {runOnNetworkError} is true.
*/
func (b eventListenerBuilder3) DisabledWhen(entityId, state string, runOnNetworkError bool) eventListenerBuilder3 {
	if entityId == "" || state == "" {
		panic(fmt.Sprintf("Either entityId or state is empty in eventListener EnabledWhen entityId='%s' state='%s' runOnNetworkError='%t'", entityId, state, runOnNetworkError))
	}
	if b.eventListener.enabledEntity != "" {
		panic(fmt.Sprintf("You can't use EnabledWhen and DisabledWhen together. Error occurred while setting DisabledWhen entityId=%s state=%s runOnNetworkError=%t", entityId, state, runOnNetworkError))
	}
	b.eventListener.disabledEntity = entityId
	b.eventListener.disabledEntityState = state
	b.eventListener.disabledEntityRunOnError = runOnNetworkError
	return b
}

func (b eventListenerBuilder3) Build() EventListener {
	return b.eventListener
}

type BaseEventMsg struct {
	Event struct {
		EventType string `json:"event_type"`
	} `json:"event"`
}

/* Functions */
func callEventListeners(app *App, msg ws.ChanMsg) {
	baseEventMsg := BaseEventMsg{}
	json.Unmarshal(msg.Raw, &baseEventMsg)
	listeners, ok := app.eventListeners[baseEventMsg.Event.EventType]
	if !ok {
		// no listeners registered for this event type
		return
	}

	for _, l := range listeners {
		// Check conditions
		if c := checkWithinTimeRange(l.betweenStart, l.betweenEnd); c.fail {
			continue
		}
		if c := checkThrottle(l.throttle, l.lastRan); c.fail {
			continue
		}
		if c := checkExceptionDates(l.exceptionDates); c.fail {
			continue
		}
		if c := checkExceptionRanges(l.exceptionRanges); c.fail {
			continue
		}
		if c := checkEnabledEntity(app.state, l.enabledEntity, l.enabledEntityState, l.enabledEntityRunOnError); c.fail {
			continue
		}
		if c := checkDisabledEntity(app.state, l.disabledEntity, l.disabledEntityState, l.disabledEntityRunOnError); c.fail {
			continue
		}

		eventData := EventData{
			Type:         baseEventMsg.Event.EventType,
			RawEventJSON: msg.Raw,
		}
		go l.callback(app.service, app.state, eventData)
		l.lastRan = carbon.Now()
	}
}
