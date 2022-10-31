package gomeassistant

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	"github.com/saml-dev/gome-assistant/internal"
	"github.com/saml-dev/gome-assistant/internal/http"
	pq "github.com/saml-dev/gome-assistant/internal/priorityqueue"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
)

type app struct {
	ctx        context.Context
	ctxCancel  context.CancelFunc
	conn       *websocket.Conn
	httpClient *http.HttpClient

	service *Service
	state   *State

	schedules         pq.PriorityQueue
	entityListeners   map[string][]*EntityListener
	entityListenersId int64
	eventListeners    map[string][]*EventListener
}

/*
DurationString represents a duration, such as "2s" or "24h".
See https://pkg.go.dev/time#ParseDuration for all valid time units.
*/
type DurationString string

/*
NewApp establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func NewApp(connString string) *app {
	token := os.Getenv("HA_AUTH_TOKEN")
	conn, ctx, ctxCancel := ws.SetupConnection(connString, token)

	httpClient := http.NewHttpClient(connString, token)

	service := NewService(conn, ctx, httpClient)
	state := newState(httpClient)

	return &app{
		conn:            conn,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		httpClient:      httpClient,
		service:         service,
		state:           state,
		schedules:       pq.New(),
		entityListeners: map[string][]*EntityListener{},
		eventListeners:  map[string][]*EventListener{},
	}
}

func (a *app) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *app) RegisterSchedule(s Schedule) {
	// realStartTime already set for sunset/sunrise
	if s.isSunrise || s.isSunset {
		a.schedules.Insert(s, float64(s.realStartTime.Unix()))
		return
	}

	if s.frequency == 0 {
		log.Fatalln("A schedule must call either Daily() or Every() when built.")
	}

	now := time.Now()
	startTime := carbon.Now().StartOfDay().Carbon2Time()
	// apply offset if set
	if s.offset.Minutes() > 0 {
		startTime = startTime.Add(s.offset)
	}

	// advance first scheduled time by frequency until it is in the future
	for startTime.Before(now) {
		startTime = startTime.Add(s.frequency)
	}

	s.realStartTime = startTime
	a.schedules.Insert(s, float64(startTime.Unix()))
}

func (a *app) RegisterEntityListener(etl EntityListener) {
	for _, entity := range etl.entityIds {
		if elList, ok := a.entityListeners[entity]; ok {
			a.entityListeners[entity] = append(elList, &etl)
		} else {
			a.entityListeners[entity] = []*EntityListener{&etl}
		}
	}
}

func (a *app) RegisterEventListener(evl EventListener) {
	for _, eventType := range evl.eventTypes {
		if elList, ok := a.eventListeners[eventType]; ok {
			a.eventListeners[eventType] = append(elList, &evl)
		} else {
			ws.SubscribeToEventType(eventType, a.conn, a.ctx)
			a.eventListeners[eventType] = []*EventListener{&evl}
		}
	}
}

func getSunriseSunset(a *app, sunrise bool, offset []DurationString) carbon.Carbon {
	printString := "Sunset"
	attrKey := "next_setting"
	if sunrise {
		printString = "Sunrise"
		attrKey = "next_rising"
	}

	var t time.Duration
	var err error
	if len(offset) == 1 {
		t, err = time.ParseDuration(string(offset[0]))
		if err != nil {
			log.Fatalf("Could not parse offset passed to %s: \"%s\"", printString, offset[0])
		}
	}

	// get next sunrise/sunset time from HA
	state, err := a.state.Get("sun.sun")
	if err != nil {
		log.Fatalln("Couldn't get sun.sun state from HA to calculate", printString)
	}

	nextSetOrRise := carbon.Parse(state.Attributes[attrKey].(string))

	// add offset if set, this code works for negative values too
	if t.Microseconds() != 0 {
		nextSetOrRise = nextSetOrRise.AddMinutes(int(t.Minutes()))
	}

	return nextSetOrRise
}

func (a *app) Start() {
	// schedules
	go runSchedules(a)

	// subscribe to state_changed events
	id := internal.GetId()
	ws.SubscribeToStateChangedEvents(id, a.conn, a.ctx)
	a.entityListenersId = id

	// entity listeners
	elChan := make(chan ws.ChanMsg)
	go ws.ListenWebsocket(a.conn, a.ctx, elChan)

	var msg ws.ChanMsg
	for {
		msg = <-elChan
		if a.entityListenersId == msg.Id {
			go callEntityListeners(a, msg.Raw)
		} else {
			go callEventListeners(a, msg)
		}
	}
}
