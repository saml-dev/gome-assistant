package gomeassistant

import (
	"context"
	"fmt"
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
	entityListeners   map[string][]*entityListener
	entityListenersId int64
	eventListeners    map[string][]*eventListener
}

type TimeString string

/*
NewApp establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func NewApp(connString string) app {
	token := os.Getenv("HA_AUTH_TOKEN")
	conn, ctx, ctxCancel := ws.SetupConnection(connString, token)

	httpClient := http.NewHttpClient(connString, token)

	service := NewService(conn, ctx, httpClient)
	state := newState(httpClient)

	return app{
		conn:            conn,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		httpClient:      httpClient,
		service:         service,
		state:           state,
		schedules:       pq.New(),
		entityListeners: map[string][]*entityListener{},
		eventListeners:  map[string][]*eventListener{},
	}
}

func (a *app) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

type ScheduleInterface interface {
	// GetNext returns the next time the schedule should execute.
	// The returned time must be in the future.
	GetNext() time.Time
	// Hash must return a string that uniquely identifies
	// a schedule.
	Hash() string
}

func (a *app) RegisterSchedule(s ScheduleInterface) {
	now := time.Now()
	startTime := s.GetNext()

	if startTime.Before(now) {
		log.Fatalln("s.GetFirst() must return time in the future")
	}

	a.schedules.Insert(s, float64(startTime.Unix()))
}

func (a *app) RegisterEntityListener(etl entityListener) {
	for _, entity := range etl.entityIds {
		if elList, ok := a.entityListeners[entity]; ok {
			a.entityListeners[entity] = append(elList, &etl)
		} else {
			a.entityListeners[entity] = []*entityListener{&etl}
		}
	}
}

func (a *app) RegisterEventListener(evl eventListener) {
	for _, eventType := range evl.eventTypes {
		if elList, ok := a.eventListeners[eventType]; ok {
			a.eventListeners[eventType] = append(elList, &evl)
		} else {
			ws.SubscribeToEventType(eventType, a.conn, a.ctx)
			a.eventListeners[eventType] = []*eventListener{&evl}
		}
	}
}

// Sunrise take an optional string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (a *app) Sunrise(offset ...TimeString) string {
	return getSunriseSunset(a, true, offset)
}

// Sunset take an optional string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (a *app) Sunset(offset ...TimeString) string {
	return getSunriseSunset(a, false, offset)
}

func getSunriseSunset(a *app, sunrise bool, offset []TimeString) string {
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

	return carbon2TimeString(nextSetOrRise)
}

func carbon2TimeString(c carbon.Carbon) string {
	return fmt.Sprintf("%02d:%02d", c.Hour(), c.Minute())
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
