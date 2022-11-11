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

type App struct {
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

type timeRange struct {
	start time.Time
	end   time.Time
}

/*
NewApp establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func NewApp(connString string) *App {
	token := os.Getenv("HA_AUTH_TOKEN")
	conn, ctx, ctxCancel := ws.SetupConnection(connString, token)

	httpClient := http.NewHttpClient(connString, token)

	service := newService(conn, ctx, httpClient)
	state := newState(httpClient)

	return &App{
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

func (a *App) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *App) RegisterSchedules(schedules ...Schedule) {
	for _, s := range schedules {
		// realStartTime already set for sunset/sunrise
		if s.isSunrise || s.isSunset {
			a.schedules.Insert(s, float64(s.realStartTime.Unix()))
			continue
		}

		if s.frequency == 0 {
			panic("A schedule must use either Daily() or Every() when built.")
		}

		now := time.Now()
		// TODO: on the day that daylight savings starts/ends,
		// this results in schedules being an hour late or hour early because
		// StartOfDay() occurs before the switch. This happens if you start
		// gome assistant on this day, and VERIFY I think will persist until you
		// start it again. Sunrise/sunset should be unaffected since those come
		// from HA.
		//
		// IDEA: splitting schedule into Interval and DailySchedule could address this on schedule
		// side, by using SetTime instead of Add(time.Duration).
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
}

func (a *App) RegisterEntityListeners(etls ...EntityListener) {
	for _, etl := range etls {
		if etl.delay != 0 && etl.toState == "" {
			panic("EntityListener error: you have to use ToState() when using Duration()")
		}

		for _, entity := range etl.entityIds {
			if elList, ok := a.entityListeners[entity]; ok {
				a.entityListeners[entity] = append(elList, &etl)
			} else {
				a.entityListeners[entity] = []*EntityListener{&etl}
			}
		}
	}
}

func (a *App) RegisterEventListeners(evls ...EventListener) {
	for _, evl := range evls {
		for _, eventType := range evl.eventTypes {
			if elList, ok := a.eventListeners[eventType]; ok {
				a.eventListeners[eventType] = append(elList, &evl)
			} else {
				ws.SubscribeToEventType(eventType, a.conn, a.ctx)
				a.eventListeners[eventType] = []*EventListener{&evl}
			}
		}
	}
}

func getSunriseSunsetFromState(s *State, sunrise bool, offset ...DurationString) carbon.Carbon {
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
			panic(fmt.Sprintf("Could not parse offset passed to %s: \"%s\"", printString, offset[0]))
		}
	}

	// get next sunrise/sunset time from HA
	state, err := s.Get("sun.sun")
	if err != nil {
		panic(fmt.Sprintf("Couldn't get sun.sun state from HA to calculate %s", printString))
	}

	nextSetOrRise := carbon.Parse(state.Attributes[attrKey].(string))

	// add offset if set, this code works for negative values too
	if t.Microseconds() != 0 {
		nextSetOrRise = nextSetOrRise.AddMinutes(int(t.Minutes()))
	}

	return nextSetOrRise
}

func getSunriseSunsetFromApp(a *App, sunrise bool, offset ...DurationString) carbon.Carbon {
	return getSunriseSunsetFromState(a.state, sunrise, offset...)
}

func (a *App) Start() {
	log.Default().Println("Starting", a.schedules.Len(), "schedules")
	log.Default().Println("Starting", len(a.entityListeners), "entity listeners")
	log.Default().Println("Starting", len(a.eventListeners), "event listeners")
	// schedules
	go runSchedules(a)

	// subscribe to state_changed events
	id := internal.GetId()
	ws.SubscribeToStateChangedEvents(id, a.conn, a.ctx)
	a.entityListenersId = id

	// entity listeners runOnStartup
	for eid, etls := range a.entityListeners {
		for _, etl := range etls {
			// ensure each ETL only runs once, even if
			// it listens to multiple entities
			if etl.runOnStartup && !etl.runOnStartupCompleted {
				entityState, err := a.state.Get(eid)
				if err != nil {
					log.Default().Println("Failed to get entity state \"", eid, "\" during startup, skipping RunOnStartup")
				}

				etl.runOnStartupCompleted = true
				go etl.callback(a.service, a.state, EntityData{
					TriggerEntityId: eid,
					FromState:       entityState.State,
					FromAttributes:  entityState.Attributes,
					ToState:         entityState.State,
					ToAttributes:    entityState.Attributes,
					LastChanged:     entityState.LastChanged,
				})
			}
		}
	}

	// entity listeners and event listeners
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
