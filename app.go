package gomeassistant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/golang-module/carbon"
	"github.com/gorilla/websocket"
	sunriseLib "github.com/nathan-osman/go-sunrise"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/http"
	ws "saml.dev/gome-assistant/internal/websocket"
)

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

	schedules         []DailySchedule
	intervals         []Interval
	entityListeners   map[string][]*EntityListener
	entityListenersId int64
	eventListeners    map[string][]*EventListener
}

// DurationString represents a duration, such as "2s" or "24h".
// See https://pkg.go.dev/time#ParseDuration for all valid time units.
type DurationString string

// TimeString is a 24-hr format time "HH:MM" such as "07:30".
type TimeString string

type timeRange struct {
	start time.Time
	end   time.Time
}

type NewAppRequest struct {
	// Required
	URL string

	// Optional
	// Deprecated: use URL instead
	// IpAddress of your Home Assistant instance i.e. "localhost"
	// or "192.168.86.59" etc.
	IpAddress string

	// Optional
	// Deprecated: use URL instead
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

// validateHomeZone verifies that the home zone entity exists and has latitude/longitude
func validateHomeZone(state State, entityID string) error {
	entity, err := state.Get(entityID)
	if err != nil {
		return fmt.Errorf("home zone entity '%s' not found: %w", entityID, err)
	}

	// Ensure it's a zone entity
	if !strings.HasPrefix(entityID, "zone.") {
		return fmt.Errorf("entity '%s' is not a zone entity (must start with zone.)", entityID)
	}

	// Verify it has latitude and longitude
	if entity.Attributes == nil {
		return fmt.Errorf("home zone entity '%s' has no attributes", entityID)
	}
	if entity.Attributes["latitude"] == nil {
		return fmt.Errorf("home zone entity '%s' missing latitude attribute", entityID)
	}
	if entity.Attributes["longitude"] == nil {
		return fmt.Errorf("home zone entity '%s' missing longitude attribute", entityID)
	}

	return nil
}

// NewApp establishes the websocket connection and returns an object
// you can use to register schedules and listeners.
func NewApp(ctx context.Context, request NewAppRequest) (*App, error) {
	if (request.URL == "" && request.IpAddress == "") || request.HAAuthToken == "" {
		slog.Error("URL and HAAuthToken are required arguments in NewAppRequest")
		return nil, ErrInvalidArgs
	}

	// Set default home zone if not provided
	if request.HomeZoneEntityId == "" {
		request.HomeZoneEntityId = "zone.home"
	}

	baseURL := &url.URL{}

	if request.URL != "" {
		var err error
		baseURL, err = url.Parse(request.URL)
		if err != nil {
			return nil, ErrInvalidArgs
		}
	} else {
		// This is deprecated and will be removed in a future release
		port := request.Port
		if port == "" {
			port = "8123"
		}
		baseURL.Scheme = "http"
		if request.Secure {
			baseURL.Scheme = "https"
		}
		baseURL.Host = request.IpAddress + ":" + port
	}

	connCtx, connCancel := context.WithTimeout(ctx, time.Second*3)
	defer connCancel()

	conn, err := ws.ConnectionFromUri(connCtx, baseURL, request.HAAuthToken)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, err
	}

	httpClient := http.NewHttpClient(baseURL, request.HAAuthToken)

	wsWriter := &ws.WebsocketWriter{Conn: conn}
	service := newService(wsWriter)
	state, err := newState(httpClient, request.HomeZoneEntityId)
	if err != nil {
		return nil, err
	}

	// Validate home zone
	if err := validateHomeZone(state, request.HomeZoneEntityId); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &App{
		conn:            conn,
		wsWriter:        wsWriter,
		ctx:             ctx,
		ctxCancel:       cancel,
		httpClient:      httpClient,
		service:         service,
		state:           state,
		entityListeners: map[string][]*EntityListener{},
		eventListeners:  map[string][]*EventListener{},
	}, nil
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
			a.schedules = append(a.schedules, s)
			continue
		}

		now := carbon.Now()
		startTime := carbon.Now().SetTimeMilli(s.hour, s.minute, 0, 0)

		// advance first scheduled time by frequency until it is in the future
		if startTime.Lt(now) {
			startTime = startTime.AddDay()
		}

		s.nextRunTime = startTime.Carbon2Time()
		a.schedules = append(a.schedules, s)
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
		a.intervals = append(a.intervals, i)
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
				ws.SubscribeToEventType(a.ctx, eventType, a.wsWriter)
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
	slog.Info("Starting", "schedules", len(a.schedules))
	slog.Info("Starting", "entity listeners", len(a.entityListeners))
	slog.Info("Starting", "event listeners", len(a.eventListeners))

	go a.runSchedules(a.ctx)
	go a.runIntervals(a.ctx)

	// subscribe to state_changed events
	id := internal.GetId()
	ws.SubscribeToStateChangedEvents(a.ctx, id, a.wsWriter)
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
