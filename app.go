package gomeassistant

import (
	"context"
	"fmt"
	"time"

	"github.com/saml-dev/gome-assistant/internal/setup"
	"nhooyr.io/websocket"
)

type app struct {
	ctx             context.Context
	ctxCancel       context.CancelFunc
	conn            *websocket.Conn
	schedules       []schedule
	entityListeners []entityListener
}

/*
App establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func App(connString string) (app, error) {
	conn, ctx, ctxCancel, err := setup.SetupConnection(connString)
	if err != nil {
		return app{}, err
	}

	return app{
		conn:            conn,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		schedules:       []schedule{},
		entityListeners: []entityListener{},
	}, nil
}

func (a *app) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *app) RegisterSchedule(s schedule) {
	fmt.Println(a.schedules)
	if s.err != nil {
		panic(s.err) // something wasn't configured properly when the schedule was built
	}

	if s.frequency == 0 {
		panic("A schedule must call either Daily() or Every() when built.")
	}

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // start at midnight today

	// apply offset if set
	if s.offset.Hour != 0 || s.offset.Minute != 0 {
		startTime.Add(time.Hour * time.Duration(s.offset.Hour))
		startTime.Add(time.Minute * time.Duration(s.offset.Minute))
	}

	// advance first scheduled time by frequency until it is in the future
	for startTime.Before(now) {
		startTime = startTime.Add(s.frequency)
	}

	s.realStartTime = startTime
	a.schedules = append(a.schedules, s)
}

const (
	FrequencyMissing time.Duration = 0

	Daily    time.Duration = time.Hour * 24
	Hourly   time.Duration = time.Hour
	Minutely time.Duration = time.Minute
)
