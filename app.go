package gomeassistant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	sunriseLib "github.com/nathan-osman/go-sunrise"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/websocket"
)

var ErrInvalidArgs = errors.New("invalid arguments provided")

// scheduledAction represents an action that can schedule and run
// itself, perhaps repeatedly.
type scheduledAction interface {
	run(ctx context.Context, app *App)
}

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	// Wraps the ws connection with added mutex locking
	conn *websocket.Conn

	httpClient *http.HttpClient

	service *Service
	state   *StateImpl

	scheduledActions   []scheduledAction
	scheduleCount      int
	entityListeners    map[string][]*EntityListener
	entitySubscription websocket.Subscription
	eventListeners     map[string][]*EventListener
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
	// EntityID of the zone representing your home e.g. "zone.home".
	// Used to pull latitude/longitude from Home Assistant
	// to calculate sunset/sunrise times.
	HomeZoneEntityID string

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
	if request.HomeZoneEntityID == "" {
		request.HomeZoneEntityID = "zone.home"
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

	conn, err := websocket.NewConn(connCtx, baseURL, request.HAAuthToken)
	if err != nil {
		return nil, err
	}

	httpClient := http.NewHttpClient(baseURL, request.HAAuthToken)

	state, err := newState(httpClient, request.HomeZoneEntityID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	app := App{
		conn:            conn,
		ctx:             ctx,
		ctxCancel:       cancel,
		httpClient:      httpClient,
		state:           state,
		entityListeners: map[string][]*EntityListener{},
		eventListeners:  map[string][]*EventListener{},
	}

	app.service = newService(&app)

	// Validate home zone
	if err := validateHomeZone(state, request.HomeZoneEntityID); err != nil {
		return nil, err
	}

	return &app, nil
}

func (app *App) Cleanup() {
	if app.ctxCancel != nil {
		app.ctxCancel()
	}
}

func (app *App) RegisterSchedules(schedules ...DailySchedule) {
	for _, s := range schedules {
		// realStartTime already set for sunset/sunrise
		if s.isSunrise || s.isSunset {
			s.nextRunTime = getNextSunRiseOrSet(app, s.isSunrise, s.sunOffset).Carbon2Time()
			app.scheduledActions = append(app.scheduledActions, s)
			app.scheduleCount++
			continue
		}

		now := carbon.Now()
		startTime := carbon.Now().SetTimeMilli(s.hour, s.minute, 0, 0)

		// advance first scheduled time by frequency until it is in the future
		if startTime.Lt(now) {
			startTime = startTime.AddDay()
		}

		s.nextRunTime = startTime.Carbon2Time()
		app.scheduledActions = append(app.scheduledActions, s)
		app.scheduleCount++
	}
}

func (app *App) RegisterIntervals(intervals ...Interval) {
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
		app.scheduledActions = append(app.scheduledActions, i)
	}
}

func (app *App) registerEntityListener(etl EntityListener) {
	if etl.delay != 0 && etl.toState == "" {
		slog.Error("EntityListener error: you have to use ToState() when using Duration()")
		panic(ErrInvalidArgs)
	}

	for _, entity := range etl.entityIDs {
		app.entityListeners[entity] = append(app.entityListeners[entity], &etl)
	}
}

func (app *App) RegisterEntityListeners(etls ...EntityListener) {
	for _, etl := range etls {
		app.registerEntityListener(etl)
	}
}

func (app *App) registerEventListener(evl EventListener) {
	for _, eventType := range evl.eventTypes {
		elList, ok := app.eventListeners[eventType]
		if !ok {
			// We're not listening to that event type yet. Ask HA to
			// send them to us, and when they arrive, call any event
			// listeners for that type (including any that are
			// registered in the future).
			eventType := eventType
			app.conn.SubscribeToEventType(
				eventType,
				func(msg websocket.ResultMessage) {
					go app.callEventListeners(eventType, msg)
				},
			)
		}
		app.eventListeners[eventType] = append(elList, &evl)
	}
}

func (app *App) RegisterEventListeners(evls ...EventListener) {
	for _, evl := range evls {
		app.registerEventListener(evl)
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

func getNextSunRiseOrSet(app *App, sunrise bool, offset ...DurationString) carbon.Carbon {
	sunriseOrSunset := getSunriseSunset(app.state, sunrise, carbon.Now(), offset...)
	if sunriseOrSunset.Lt(carbon.Now()) {
		// if we're past today's sunset or sunrise (accounting for offset) then get tomorrows
		// as that's the next time the schedule will run
		sunriseOrSunset = getSunriseSunset(app.state, sunrise, carbon.Tomorrow(), offset...)
	}
	return sunriseOrSunset
}

func (app *App) Start() {
	slog.Info("Starting", "schedules", app.scheduleCount)
	slog.Info("Starting", "entity listeners", len(app.entityListeners))
	slog.Info("Starting", "event listeners", len(app.eventListeners))

	go app.runScheduledActions(app.ctx)

	// subscribe to state_changed events
	app.entitySubscription = app.conn.SubscribeToStateChangedEvents(
		func(msg websocket.ResultMessage) {
			go app.callEntityListeners(msg.Raw)
		},
	)

	// entity listeners runOnStartup
	for eid, etls := range app.entityListeners {
		for _, etl := range etls {
			// ensure each ETL only runs once, even if
			// it listens to multiple entities
			if etl.runOnStartup && !etl.runOnStartupCompleted {
				entityState, err := app.state.Get(eid)
				if err != nil {
					slog.Warn("Failed to get entity state \"", eid, "\" during startup, skipping RunOnStartup")
				}

				etl.runOnStartupCompleted = true
				go etl.callback(app.service, app.state, EntityData{
					TriggerEntityID: eid,
					FromState:       entityState.State,
					FromAttributes:  entityState.Attributes,
					ToState:         entityState.State,
					ToAttributes:    entityState.Attributes,
					LastChanged:     entityState.LastChanged,
				})
			}
		}
	}

	// Start listen on the connection for incoming messages:
	if err := app.conn.Run(); err != nil {
		slog.Error("Error reading from websocket", "err", err)
	}
}

// runScheduledActions starts a goroutine to run each `DailySchedule`
// and each `Interval` that has been configured. The `run()` method of
// each of those instances takes care of deciding when to run and
// invoking its callback.
func (app *App) runScheduledActions(ctx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, action := range app.scheduledActions {
		wg.Add(1)
		go func(action scheduledAction) {
			defer wg.Done()
			action.run(ctx, app)
		}(action)
	}
}

func (app *App) GetService() *Service {
	return app.service
}

func (app *App) GetState() State {
	return app.state
}
