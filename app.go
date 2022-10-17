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
	entityListeners   map[string][]entityListener
	entityListenersId int64
}

/*
Time is a 24-hr format string with hour and minute,
e.g. "07:00" for 7AM or "23:00" for 11PM.
*/
type Time string

/*
NewApp establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func NewApp(connString string) app {
	token := os.Getenv("HA_AUTH_TOKEN")
	conn, ctx, ctxCancel := ws.SetupConnection(connString, token)

	httpClient := http.NewHttpClient(connString, token)

	service := NewService(conn, ctx, httpClient)
	state := NewState(httpClient)

	return app{
		conn:            conn,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		httpClient:      httpClient,
		service:         service,
		state:           state,
		schedules:       pq.New(),
		entityListeners: map[string][]entityListener{},
	}
}

func (a *app) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *app) RegisterSchedule(s schedule) {
	if s.err != nil {
		log.Fatalln(s.err) // something wasn't configured properly when the schedule was built
	}

	if s.frequency == 0 {
		log.Fatalln("A schedule must call either Daily() or Every() when built.")
	}

	// TODO: consider moving all time stuff to carbon?
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // start at midnight today
	// apply offset if set
	if s.offset.Minutes() > 0 {
		startTime = startTime.Add(s.offset)
	}

	// advance first scheduled time by frequency until it is in the future
	for startTime.Before(now) {
		startTime = startTime.Add(s.frequency)
	}

	s.realStartTime = startTime
	a.schedules.Insert(s, float64(startTime.Unix())) // TODO: this blows up because schedule can't be used as key for map in prio queue lib. Just copy/paste and tweak as needed
}

func (a *app) RegisterEntityListener(el entityListener) {
	for _, entity := range el.entityIds {
		if elList, ok := a.entityListeners[entity]; ok {
			a.entityListeners[entity] = append(elList, el)
		} else {
			a.entityListeners[entity] = []entityListener{el}
		}
	}
}

// Sunrise take an optional string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (a *app) Sunrise(offset ...string) Time {
	return getSunriseSunset(a, true, offset)
}

// Sunset take an optional string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (a *app) Sunset(offset ...string) Time {
	return getSunriseSunset(a, false, offset)
}

func getSunriseSunset(a *app, sunrise bool, offset []string) Time {
	printString := "Sunset"
	attrKey := "next_setting"
	if sunrise {
		printString = "Sunrise"
		attrKey = "next_rising"
	}
	var t time.Duration
	var err error
	if len(offset) == 1 {
		t, err = time.ParseDuration(offset[0])
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
	log.Default().Println(nextSetOrRise)

	// add offset if set, this code works for negative values too
	if t.Microseconds() != 0 {
		nextSetOrRise = nextSetOrRise.AddMinutes(int(t.Minutes()))
		log.Default().Println(nextSetOrRise)
	}

	return carbon2TimeString(nextSetOrRise)
}

func carbon2TimeString(c carbon.Carbon) Time {
	return Time(fmt.Sprintf("%02d:%02d", c.Hour(), c.Minute()))
}

type subEvent struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

func (a *app) Start() {
	// schedules
	go RunSchedules(a)

	// subscribe to state_changed events
	id := internal.GetId()
	e := subEvent{
		Id:        id,
		Type:      "subscribe_events",
		EventType: "state_changed",
	}
	ws.WriteMessage(e, a.conn, a.ctx)
	a.entityListenersId = id

	// entity listeners
	elChan := make(chan ws.ChanMsg)
	go ws.ListenWebsocket(a.conn, a.ctx, elChan)

	var msg ws.ChanMsg
	for {
		msg = <-elChan
		if a.entityListenersId == msg.Id {
			go callEntityListeners(a, msg.Raw)
		}
	}

	// NOTE:should the prio queue and websocket listener both write to a channel or something?
	// then select from that and spawn new goroutine to call callback?

	// TODO: loop through schedules and create heap priority queue

	// TODO: figure out looping listening to messages for
	// listeners
}

const (
	FrequencyMissing time.Duration = 0

	Daily    time.Duration = time.Hour * 24
	Hourly   time.Duration = time.Hour
	Minutely time.Duration = time.Minute
)
