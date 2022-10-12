package gomeassistant

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/saml-dev/gome-assistant/internal/http"
	ws "github.com/saml-dev/gome-assistant/internal/websocket"
	"nhooyr.io/websocket"
)

type app struct {
	ctx        context.Context
	ctxCancel  context.CancelFunc
	conn       *websocket.Conn
	httpClient *http.HttpClient

	service *Service
	state   *State

	schedules       []schedule
	entityListeners []entityListener
}

/*
App establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func App(connString string) app {
	token := os.Getenv("AUTH_TOKEN")
	conn, ctx, ctxCancel := ws.SetupConnection(connString)

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
		schedules:       []schedule{},
		entityListeners: []entityListener{},
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
		startTime.Add(s.offset)
	}

	// advance first scheduled time by frequency until it is in the future
	for startTime.Before(now) {
		startTime = startTime.Add(s.frequency)
	}

	s.realStartTime = startTime
	a.schedules = append(a.schedules, s)
}

func (a *app) Start() {
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
