package gomeassistant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	sunriseLib "github.com/nathan-osman/go-sunrise"
	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/http"
	pq "saml.dev/gome-assistant/internal/priorityqueue"
	ws "saml.dev/gome-assistant/internal/websocket"
)

// Returned by NewApp() if authentication fails
var ErrInvalidToken = ws.ErrInvalidToken

var ErrInvalidArgs = errors.New("invalid arguments provided")

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *websocket.Conn

	// Wraps the ws connection with added mutex locking
	wsWriter *ws.WebsocketWriter

	httpClient *http.HttpClient

	service *Service
	state   *StateImpl

	schedules         pq.PriorityQueue
	intervals         pq.PriorityQueue
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
TimeString is a 24-hr format time "HH:MM" such as "07:30".
*/
type TimeString string

type timeRange struct {
	start time.Time
	end   time.Time
}

type NewAppConfig struct {
	// RESTBaseURI is the base URI for REST requests; for example,
	//  * `http://homeassistant.local:8123/api` from outside of the
	//    HA appliance (without encryption)
	//  * `https://homeassistant.local:8123/api` from outside of the
	//    HA appliance (with encryption)
	//  * `http://supervisor/core/api` from an add-on running within
	//    the appliance and connecting via the proxy
	RESTBaseURI string

	// WebsocketURI is the base URI for websocket connections; for
	// example,
	//  * `ws://homeassistant.local:8123/api/websocket` from outside
	//    of the HA appliance (without encryption)
	//  * `wss://homeassistant.local:8123/api/websocket` from outside
	//    of the HA appliance (with encryption)
	//  * `ws://supervisor/core/api/websocket` from an add-on running
	//    within the appliance and connecting via the proxy
	WebsocketURI string

	// Auth token generated in Home Assistant. Used to connect to the
	// Websocket API.
	HAAuthToken string

	// Required
	// EntityId of the zone representing your home e.g. "zone.home".
	// Used to pull latitude/longitude from Home Assistant
	// to calculate sunset/sunrise times.
	HomeZoneEntityId string
}

/*
NewAppFromConfig establishes the websocket connection and returns an
object you can use to register schedules and listeners, based on the
URIs that it should connect to.
*/
func NewAppFromConfig(config NewAppConfig) (*App, error) {
	if config.RESTBaseURI == "" || config.WebsocketURI == "" ||
		config.HAAuthToken == "" || config.HomeZoneEntityId == "" {
		slog.Error("RESTBaseURI, WebsocketURI, HAAuthToken, and HomeZoneEntityId are all required arguments in NewAppRequest")
		return nil, ErrInvalidArgs
	}

	conn, ctx, ctxCancel, err := ws.ConnectionFromUri(config.WebsocketURI, config.HAAuthToken)
	if err != nil {
		return nil, err
	}

	httpClient := http.ClientFromUri(config.RESTBaseURI, config.HAAuthToken)

	wsWriter := &ws.WebsocketWriter{Conn: conn}
	service := newService(wsWriter, httpClient)
	state, err := newState(httpClient, config.HomeZoneEntityId)
	if err != nil {
		return nil, err
	}

	return &App{
		conn:            conn,
		wsWriter:        wsWriter,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		httpClient:      httpClient,
		service:         service,
		state:           state,
		schedules:       pq.New(),
		intervals:       pq.New(),
		entityListeners: map[string][]*EntityListener{},
		eventListeners:  map[string][]*EventListener{},
	}, nil
}

type NewAppRequest struct {
	// Required
	// IpAddress of your Home Assistant instance i.e. "localhost"
	// or "192.168.86.59" etc.
	IpAddress string

	// Optional
	// Port number Home Assistant is running on. Defaults to 8123.
	Port string

	// Required
	// Auth token generated in Home Assistant. Used
	// to connect to the Websocket API.
	HAAuthToken string

	// Required
	// EntityId of the zone representing your home e.g. "zone.home".
	// Used to pull latitude/longitude from Home Assistant
	// to calculate sunset/sunrise times.
	HomeZoneEntityId string

	// Optional
	// Whether to use secure connections for http and websockets.
	// Setting this to `true` will use `https://` instead of `https://`
	// and `wss://` instead of `ws://`.
	Secure bool
}

/*
NewApp establishes the websocket connection and returns an object
you can use to register schedules and listeners.
*/
func NewApp(request NewAppRequest) (*App, error) {
	if request.IpAddress == "" || request.HAAuthToken == "" || request.HomeZoneEntityId == "" {
		slog.Error("IpAddress, HAAuthToken, and HomeZoneEntityId are all required arguments in NewAppRequest")
		return nil, ErrInvalidArgs
	}
	port := request.Port
	if port == "" {
		port = "8123"
	}

	config := NewAppConfig{
		HAAuthToken:      request.HAAuthToken,
		HomeZoneEntityId: request.HomeZoneEntityId,
	}

	if request.Secure {
		config.WebsocketURI = fmt.Sprintf("wss://%s:%s/api/websocket", request.IpAddress, port)
		config.RESTBaseURI = fmt.Sprintf("https://%s:%s/api", request.IpAddress, port)
	} else {
		config.WebsocketURI = fmt.Sprintf("ws://%s:%s/api/websocket", request.IpAddress, port)
		config.RESTBaseURI = fmt.Sprintf("http://%s:%s/api", request.IpAddress, port)
	}

	return NewAppFromConfig(config)
}

func (a *App) Cleanup() {
	if a.ctxCancel != nil {
		a.ctxCancel()
	}
}

func (a *App) RegisterSchedules(schedules ...DailySchedule) {
	for _, s := range schedules {
		// realStartTime already set for sunset/sunrise
		if s.isSunrise || s.isSunset {
			s.nextRunTime = getNextSunRiseOrSet(a, s.isSunrise, s.sunOffset).Carbon2Time()
			a.schedules.Insert(s, float64(s.nextRunTime.Unix()))
			continue
		}

		now := carbon.Now()
		startTime := carbon.Now().SetTimeMilli(s.hour, s.minute, 0, 0)

		// advance first scheduled time by frequency until it is in the future
		if startTime.Lt(now) {
			startTime = startTime.AddDay()
		}

		s.nextRunTime = startTime.Carbon2Time()
		a.schedules.Insert(s, float64(startTime.Carbon2Time().Unix()))
	}
}

func (a *App) RegisterIntervals(intervals ...Interval) {
	for _, i := range intervals {
		if i.frequency == 0 {
			slog.Error("A schedule must use either set frequency via Every()")
			panic(ErrInvalidArgs)
		}

		i.nextRunTime = internal.ParseTime(string(i.startTime)).Carbon2Time()
		now := time.Now()
		for i.nextRunTime.Before(now) {
			i.nextRunTime = i.nextRunTime.Add(i.frequency)
		}
		a.intervals.Insert(i, float64(i.nextRunTime.Unix()))
	}
}

func (a *App) RegisterEntityListeners(etls ...EntityListener) {
	for _, etl := range etls {
		etl := etl
		if etl.delay != 0 && etl.toState == "" {
			slog.Error("EntityListener error: you have to use ToState() when using Duration()")
			panic(ErrInvalidArgs)
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
		evl := evl
		for _, eventType := range evl.eventTypes {
			if elList, ok := a.eventListeners[eventType]; ok {
				a.eventListeners[eventType] = append(elList, &evl)
			} else {
				ws.SubscribeToEventType(eventType, a.wsWriter)
				a.eventListeners[eventType] = []*EventListener{&evl}
			}
		}
	}
}

func getSunriseSunset(s *StateImpl, sunrise bool, dateToUse carbon.Carbon, offset ...DurationString) carbon.Carbon {
	date := dateToUse.Carbon2Time()
	rise, set := sunriseLib.SunriseSunset(s.latitude, s.longitude, date.Year(), date.Month(), date.Day())
	rise, set = rise.Local(), set.Local()

	val := set
	printString := "Sunset"
	if sunrise {
		val = rise
		printString = "Sunrise"
	}

	setOrRiseToday := carbon.Parse(val.String())

	var t time.Duration
	var err error
	if len(offset) == 1 {
		t, err = time.ParseDuration(string(offset[0]))
		if err != nil {
			parsingErr := fmt.Errorf("could not parse offset passed to %s: \"%s\": %w", printString, offset[0], err)
			slog.Error(parsingErr.Error())
			panic(parsingErr)
		}
	}

	// add offset if set, this code works for negative values too
	if t.Microseconds() != 0 {
		setOrRiseToday = setOrRiseToday.AddMinutes(int(t.Minutes()))
	}

	return setOrRiseToday
}

func getNextSunRiseOrSet(a *App, sunrise bool, offset ...DurationString) carbon.Carbon {
	sunriseOrSunset := getSunriseSunset(a.state, sunrise, carbon.Now(), offset...)
	if sunriseOrSunset.Lt(carbon.Now()) {
		// if we're past today's sunset or sunrise (accounting for offset) then get tomorrows
		// as that's the next time the schedule will run
		sunriseOrSunset = getSunriseSunset(a.state, sunrise, carbon.Tomorrow(), offset...)
	}
	return sunriseOrSunset
}

func (a *App) Start() {
	slog.Info("Starting", "schedules", a.schedules.Len())
	slog.Info("Starting", "entity listeners", len(a.entityListeners))
	slog.Info("Starting", "event listeners", len(a.eventListeners))

	go runSchedules(a)
	go runIntervals(a)

	// subscribe to state_changed events
	id := internal.GetId()
	ws.SubscribeToStateChangedEvents(id, a.wsWriter)
	a.entityListenersId = id

	// entity listeners runOnStartup
	for eid, etls := range a.entityListeners {
		for _, etl := range etls {
			// ensure each ETL only runs once, even if
			// it listens to multiple entities
			if etl.runOnStartup && !etl.runOnStartupCompleted {
				entityState, err := a.state.Get(eid)
				if err != nil {
					slog.Warn("Failed to get entity state \"", eid, "\" during startup, skipping RunOnStartup")
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
	go ws.ListenWebsocket(a.conn, elChan)

	for {
		msg, ok := <-elChan
		if !ok {
			break
		}
		if a.entityListenersId == msg.Id {
			go callEntityListeners(a, msg.Raw)
		} else {
			go callEventListeners(a, msg)
		}
	}
}

func (a *App) GetService() *Service {
	return a.service
}

func (a *App) GetState() State {
	return a.state
}
