package gomeassistant

// TODO: impl eventListener. could probably create generic listener struct for
// code reuse between eventListener and eventListener

import (
	"encoding/json"
	"log"
	"time"

	"github.com/golang-module/carbon"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

type eventListener struct {
	eventTypes   []string
	callback     eventListenerCallback
	betweenStart string
	betweenEnd   string
	throttle     time.Duration
	lastRan      carbon.Carbon
}

type eventListenerCallback func(*Service, EventData)

type EventData struct {
	Type         string
	RawEventJSON []byte
}

/* Methods */

func EventListenerBuilder() eventListenerBuilder1 {
	return eventListenerBuilder1{eventListener{
		lastRan: carbon.Now().StartOfCentury(),
	}}
}

type eventListenerBuilder1 struct {
	eventListener
}

func (b eventListenerBuilder1) EventTypes(ets ...string) eventListenerBuilder2 {
	b.eventTypes = ets
	return eventListenerBuilder2(b)
}

type eventListenerBuilder2 struct {
	eventListener
}

func (b eventListenerBuilder2) Call(callback eventListenerCallback) eventListenerBuilder3 {
	b.eventListener.callback = callback
	return eventListenerBuilder3(b)
}

type eventListenerBuilder3 struct {
	eventListener
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

func (b eventListenerBuilder3) Throttle(s TimeString) eventListenerBuilder3 {
	d, err := time.ParseDuration(string(s))
	if err != nil {
		log.Fatalf("Couldn't parse string duration passed to Throttle(): \"%s\" see https://pkg.go.dev/time#ParseDuration for valid time units", s)
	}
	b.eventListener.throttle = d
	return b
}

func (b eventListenerBuilder3) Build() eventListener {
	return b.eventListener
}

type BaseEventMsg struct {
	Event struct {
		EventType string `json:"event_type"`
	} `json:"event"`
}

/* Functions */
func callEventListeners(app *app, msg ws.ChanMsg) {
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
			return
		}
		if c := checkThrottle(l.throttle, l.lastRan); c.fail {
			return
		}

		eventData := EventData{
			Type:         baseEventMsg.Event.EventType,
			RawEventJSON: msg.Raw,
		}
		go l.callback(app.service, eventData)
		l.lastRan = carbon.Now()
	}
}
