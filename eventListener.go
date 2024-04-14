package gomeassistant

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-module/carbon"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/websocket"
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

	enabledEntities  []internal.EnabledDisabledInfo
	disabledEntities []internal.EnabledDisabledInfo
}

type EventListenerCallback func(*Service, State, EventData)

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

func (b eventListenerBuilder3) ExceptionDates(
	t time.Time, tl ...time.Time,
) eventListenerBuilder3 {
	b.eventListener.exceptionDates = append(tl, t)
	return b
}

func (b eventListenerBuilder3) ExceptionRange(start, end time.Time) eventListenerBuilder3 {
	b.eventListener.exceptionRanges = append(
		b.eventListener.exceptionRanges,
		timeRange{start, end},
	)
	return b
}

// Enable this listener only when the current state of {entityID}
// matches {state}. If there is a network error while retrieving
// state, the listener runs if {runOnNetworkError} is true.
func (b eventListenerBuilder3) EnabledWhen(
	entityID, state string, runOnNetworkError bool,
) eventListenerBuilder3 {
	if entityID == "" {
		panic(
			fmt.Sprintf(
				"entityID is empty in eventListener EnabledWhen "+
					"entityID='%s' state='%s' runOnNetworkError='%t'",
				entityID, state, runOnNetworkError,
			),
		)
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.eventListener.enabledEntities = append(b.eventListener.enabledEntities, i)
	return b
}

// Disable this listener when the current state of {entityID} matches
// {state}. If there is a network error while retrieving state, the
// listener runs if {runOnNetworkError} is true.
func (b eventListenerBuilder3) DisabledWhen(
	entityID, state string, runOnNetworkError bool,
) eventListenerBuilder3 {
	if entityID == "" {
		panic(
			fmt.Sprintf(
				"entityID is empty in eventListener EnabledWhen "+
					"entityID='%s' state='%s' runOnNetworkError='%t'",
				entityID, state, runOnNetworkError,
			),
		)
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.eventListener.disabledEntities = append(b.eventListener.disabledEntities, i)
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
func (app *App) callEventListeners(msg websocket.ChanMsg) {
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
		if c := checkEnabledEntity(app.state, l.enabledEntities); c.fail {
			continue
		}
		if c := checkDisabledEntity(app.state, l.disabledEntities); c.fail {
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
