package gomeassistant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-module/carbon"
	sunriseLib "github.com/nathan-osman/go-sunrise"

	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/websocket"
)

var ErrInvalidArgs = errors.New("invalid arguments provided")

// ErrAppClosed is returned when Start is called after Cleanup, or when Cleanup
// was called before the application had a chance to start.
var ErrAppClosed = errors.New("app is closed")

// ErrAppNotRunning is returned when an operation requires an active session
// but Start has not acquired a connection (or has already released it).
var ErrAppNotRunning = errors.New("app is not running")

// ErrConnectionClosed is returned when a connection terminates without
// recording a more specific terminal error.
var ErrConnectionClosed = errors.New("websocket connection is closed")

// scheduledAction represents an action that can schedule and run
// itself, perhaps repeatedly.
type scheduledAction interface {
	run(ctx context.Context, app *App)
}

type App struct {
	closed  atomic.Bool
	workers sync.WaitGroup

	ctx context.Context
	run atomic.Pointer[runCancellation]

	baseURL          *url.URL
	authToken        string
	homeZoneEntityID string

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

// runCancellation is published before Start begins startup work so Cleanup can
// always cancel a Start call that has passed its initial closed check.
type runCancellation struct {
	cancel context.CancelFunc
}

// activeConn returns the connection owned by the current Start call.
func (app *App) activeConn() (*websocket.Conn, error) {
	if app.closed.Load() {
		return nil, ErrAppClosed
	}
	if app.conn == nil {
		return nil, ErrAppNotRunning
	}
	return app.conn, nil
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

	// Required
	// Auth token generated in Home Assistant. Used
	// to connect to the Websocket API.
	HAAuthToken string

	// Required
	// EntityID of the zone representing your home e.g. "zone.home".
	// Used to pull latitude/longitude from Home Assistant
	// to calculate sunset/sunrise times.
	HomeZoneEntityID string
}

// NewApp validates its configuration and returns an inert application that
// you can use to register schedules and listeners.
func NewApp(ctx context.Context, request NewAppRequest) (*App, error) {
	if ctx == nil || request.URL == "" || request.HAAuthToken == "" {
		slog.Error("URL and HAAuthToken are required arguments in NewAppRequest")
		return nil, ErrInvalidArgs
	}

	// Set default home zone if not provided
	if request.HomeZoneEntityID == "" {
		request.HomeZoneEntityID = "zone.home"
	}

	baseURL, err := url.Parse(request.URL)
	if err != nil {
		return nil, ErrInvalidArgs
	}

	if (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		return nil, ErrInvalidArgs
	}

	httpClient := http.NewHttpClient(baseURL, request.HAAuthToken)
	state := newState(httpClient)

	app := App{
		baseURL:          baseURL,
		authToken:        request.HAAuthToken,
		homeZoneEntityID: request.HomeZoneEntityID,
		httpClient:       httpClient,
		state:            state,
		entityListeners:  map[string][]*EntityListener{},
		eventListeners:   map[string][]*EventListener{},
	}

	app.service = newService(&app)

	return &app, nil
}

// Cleanup permanently closes the app and asks the current Start call to stop.
// Start owns connection cleanup and waits for application-owned goroutines
// before returning.
func (app *App) Cleanup() {
	if app.closed.Swap(true) {
		return
	}
	if run := app.run.Load(); run != nil {
		run.cancel()
	}
}

func (app *App) RegisterSchedules(schedules ...DailySchedule) {
	for _, s := range schedules {
		// Keep scheduler state internal: registrations take values, and workers
		// advance this dedicated copy across Start sessions.
		schedule := new(DailySchedule)
		*schedule = s

		// Solar schedules depend on the home-zone coordinates loaded by Start.
		// Keep registration network-free and initialize them after that load.
		if schedule.isSunrise || schedule.isSunset {
			app.scheduledActions = append(app.scheduledActions, schedule)
			app.scheduleCount++
			continue
		}

		now := carbon.Now()
		startTime := carbon.Now().SetTimeMilli(schedule.hour, schedule.minute, 0, 0)

		// advance first scheduled time by frequency until it is in the future
		if startTime.Lt(now) {
			startTime = startTime.AddDay()
		}

		schedule.nextRunTime = startTime.Carbon2Time()
		app.scheduledActions = append(app.scheduledActions, schedule)
		app.scheduleCount++
	}
}

func (app *App) RegisterIntervals(intervals ...Interval) {
	for _, i := range intervals {
		// Keep scheduler state internal: registrations take values, and workers
		// advance this dedicated copy across Start sessions.
		interval := new(Interval)
		*interval = i

		if interval.frequency == 0 {
			slog.Error("A schedule must use either set frequency via Every()")
			panic(ErrInvalidArgs)
		}

		interval.nextRunTime = internal.ParseTime(string(interval.startTime)).Carbon2Time()
		now := time.Now()
		for interval.nextRunTime.Before(now) {
			interval.nextRunTime = interval.nextRunTime.Add(interval.frequency)
		}
		app.scheduledActions = append(app.scheduledActions, interval)
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
		app.eventListeners[eventType] = append(app.eventListeners[eventType], &evl)
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

func (app *App) initializeSolarSchedules() {
	for _, action := range app.scheduledActions {
		schedule, ok := action.(*DailySchedule)
		if !ok || (!schedule.isSunrise && !schedule.isSunset) {
			continue
		}
		schedule.nextRunTime = getNextSunRiseOrSet(app, schedule.isSunrise, schedule.sunOffset).Carbon2Time()
	}
}

// Start owns one Home Assistant connection and blocks until the context is
// canceled or the session ends. It does not retry; callers decide whether and
// when to call Start again. Start calls must not overlap.
func (app *App) Start(ctx context.Context) error {
	if ctx == nil {
		return ErrInvalidArgs
	}

	if app.closed.Load() {
		return ErrAppClosed
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	runCtx, cancel := context.WithCancel(ctx)
	run := &runCancellation{cancel: cancel}
	app.run.Store(run)
	if app.closed.Load() {
		cancel()
	}
	app.ctx = runCtx

	defer func() {
		cancel()
		if app.conn != nil {
			_ = app.conn.Close()
		}
		app.workers.Wait()
		app.conn = nil
		app.run.CompareAndSwap(run, nil)
	}()

	slog.Info("Starting", "schedules", app.scheduleCount)
	slog.Info("Starting", "entity listeners", len(app.entityListeners))
	slog.Info("Starting", "event listeners", len(app.eventListeners))

	conn, err := websocket.NewConn(runCtx, app.baseURL, app.authToken)
	if err != nil {
		if runErr := runCtx.Err(); runErr != nil {
			return runErr
		}
		return err
	}
	app.conn = conn

	if err := app.state.loadHomeZone(runCtx, app.homeZoneEntityID); err != nil {
		if runErr := runCtx.Err(); runErr != nil {
			return runErr
		}
		return err
	}
	app.initializeSolarSchedules()
	if err := app.subscribeSession(conn); err != nil {
		if runErr := runCtx.Err(); runErr != nil {
			return runErr
		}
		return err
	}
	app.goTracked(func() { app.runScheduledActions(runCtx) })
	app.runStartupCallbacks(runCtx)
	return conn.Run(runCtx)
}

func (app *App) subscribeSession(conn *websocket.Conn) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("subscribe during startup: %v", recovered)
		}
	}()

	// subscribe to state_changed events
	app.entitySubscription = conn.SubscribeToStateChangedEvents(
		func(msg websocket.Message) {
			app.goTracked(func() { app.callEntityListeners(msg.Raw) })
		},
	)

	eventTypes := make([]string, 0, len(app.eventListeners))
	for eventType := range app.eventListeners {
		eventTypes = append(eventTypes, eventType)
	}
	for _, eventType := range eventTypes {
		eventType := eventType
		conn.SubscribeToEventType(eventType, func(msg websocket.Message) {
			if msg.Type != "event" {
				return
			}
			app.goTracked(func() { app.callEventListeners(eventType, msg) })
		})
	}
	return nil
}

func (app *App) runStartupCallbacks(ctx context.Context) {
	completed := make(map[*EntityListener]bool)

	// entity listeners runOnStartup
	for eid, etls := range app.entityListeners {
		for _, etl := range etls {
			// ensure each ETL only runs once, even if
			// it listens to multiple entities
			if etl.runOnStartup && !completed[etl] {
				entityState, err := app.state.getWithContext(ctx, eid)
				if err != nil {
					slog.Warn("Failed to get entity state \"", eid, "\" during startup, skipping RunOnStartup")
					continue
				}

				completed[etl] = true
				data := EntityData{
					TriggerEntityID: eid,
					FromState:       entityState.State,
					FromAttributes:  entityState.Attributes,
					ToState:         entityState.State,
					ToAttributes:    entityState.Attributes,
					LastChanged:     entityState.LastChanged,
				}
				app.goTracked(func() { etl.callback(app.service, app.state, data) })
			}
		}
	}
}

func (app *App) goTracked(fn func()) {
	if app.ctx != nil && app.ctx.Err() != nil {
		return
	}
	app.workers.Add(1)
	go func() {
		defer app.workers.Done()
		fn()
	}()
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
