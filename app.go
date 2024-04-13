package gomeassistant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	sunriseLib "github.com/nathan-osman/go-sunrise"
	"golang.org/x/sync/errgroup"
	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/internal/http"
	pq "saml.dev/gome-assistant/internal/priorityqueue"
	ws "saml.dev/gome-assistant/internal/websocket"
)

// Returned by NewApp() if authentication fails
var ErrInvalidToken = ws.ErrInvalidToken

var ErrInvalidArgs = errors.New("invalid arguments provided")

type App struct {
	// Wraps the ws connection with added mutex locking
	wsConn *ws.WebsocketConn

	httpClient *http.HttpClient

	service *Service
	state   *StateImpl

	scheduledActions  pq.PriorityQueue
	entityListeners   map[string][]*EntityListener
	entityListenersId int64
	eventListeners    map[string][]*EventListener

	// If `App.Start()` has been called, `cancel()` cancels the
	// context being used, which causes the app to shut down cleanly.
	cancel context.CancelFunc

	closeOnce sync.Once
}

// DurationString represents a duration, such as "2s" or "24h". See
// https://pkg.go.dev/time#ParseDuration for all valid time units.
type DurationString string

// TimeString is a 24-hr format time "HH:MM" such as "07:30".
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

// NewAppFromConfig establishes the websocket connection and returns
// an object you can use to register schedules and listeners, based on
// the URIs that it should connect to. `ctx` is used only to limit the
// time spent connecting; it cannot be used after that to cancel the
// app.
func NewAppFromConfig(ctx context.Context, config NewAppConfig) (*App, error) {
	if config.RESTBaseURI == "" || config.WebsocketURI == "" ||
		config.HAAuthToken == "" || config.HomeZoneEntityId == "" {
		slog.Error("RESTBaseURI, WebsocketURI, HAAuthToken, and HomeZoneEntityId are all required arguments in NewAppRequest")
		return nil, ErrInvalidArgs
	}

	conn, err := ws.ConnectionFromUri(ctx, config.WebsocketURI, config.HAAuthToken)
	if err != nil {
		return nil, err
	}

	httpClient := http.ClientFromUri(config.RESTBaseURI, config.HAAuthToken)

	wsWriter := &ws.WebsocketConn{Conn: conn}
	service := newService(wsWriter, httpClient)
	state, err := newState(httpClient, config.HomeZoneEntityId)
	if err != nil {
		return nil, err
	}

	return &App{
		wsConn:           wsWriter,
		httpClient:       httpClient,
		service:          service,
		state:            state,
		scheduledActions: pq.New(),
		entityListeners:  map[string][]*EntityListener{},
		eventListeners:   map[string][]*EventListener{},
		cancel:           func() {},
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

// NewApp establishes the websocket connection and returns an object
// you can use to register schedules and listeners. `ctx` is used only
// to limit the time spent connecting; it cannot be used after that to
// cancel the app. If this function returns successfully, then
// `App.Close()` must eventually be called to release resources.
func NewApp(ctx context.Context, request NewAppRequest) (*App, error) {
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

	return NewAppFromConfig(ctx, config)
}

type scheduledAction interface {
	String() string
	Hash() string
	initializeNextRunTime(a *App)
	shouldRun(a *App) bool
	run(a *App)
	updateNextRunTime(a *App)
	getNextRunTime() time.Time
}

func (a *App) RegisterScheduledAction(action scheduledAction) {
	action.initializeNextRunTime(a)
	a.scheduledActions.Insert(action, float64(action.getNextRunTime().Unix()))
}

func (a *App) RegisterSchedules(schedules ...*DailySchedule) {
	for _, s := range schedules {
		a.RegisterScheduledAction(s)
	}
}

func (a *App) RegisterIntervals(intervals ...*Interval) {
	for _, i := range intervals {
		a.RegisterScheduledAction(i)
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
				ws.SubscribeToEventType(eventType, a.wsConn)
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

// Start the app. When `ctx` expires, the app closes the connection
// and returns.
func (a *App) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	a.cancel = cancel
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	slog.Info("Starting", "scheduled actions", a.scheduledActions.Len())
	slog.Info("Starting", "entity listeners", len(a.entityListeners))
	slog.Info("Starting", "event listeners", len(a.eventListeners))

	eg.Go(func() error {
		a.runScheduledActions(ctx)
		return nil
	})

	// subscribe to state_changed events
	id := internal.GetId()
	ws.SubscribeToStateChangedEvents(id, a.wsConn)
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
				etl := etl
				eg.Go(func() error {
					etl.callback(a.service, a.state, EntityData{
						TriggerEntityId: eid,
						FromState:       entityState.State,
						FromAttributes:  entityState.Attributes,
						ToState:         entityState.State,
						ToAttributes:    entityState.Attributes,
						LastChanged:     entityState.LastChanged,
					})
					return nil
				})
			}
		}
	}

	// entity listeners and event listeners
	elChan := make(chan ws.ChanMsg)
	eg.Go(func() error {
		a.wsConn.ListenWebsocket(elChan)
		cancel()
		return nil
	})

	eg.Go(func() error {
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
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		a.Close()
		return nil
	})

	eg.Wait()
}

// Close closes the connection and releases any resources. It may be
// called more than once; only the first call does anything.
func (a *App) Close() {
	a.closeOnce.Do(func() {
		a.close()
	})
}

// close closes the connection and releases resources. It must be
// called exactly once.
func (a *App) close() {
	a.cancel()
	a.wsConn.Close()
}

func (a *App) GetService() *Service {
	return a.service
}

func (a *App) GetState() State {
	return a.state
}

func (a *App) runScheduledActions(ctx context.Context) {
	if a.scheduledActions.Len() == 0 {
		return
	}

	// Create a new, but stopped, timer:
	timer := time.NewTimer(1 * time.Hour)
	if !timer.Stop() {
		<-timer.C
	}

	for {
		action := a.popScheduledAction()
		if action.getNextRunTime().After(time.Now()) {
			timer.Reset(time.Until(action.getNextRunTime()))

			select {
			case <-timer.C:
			case <-ctx.Done():
				return
			}
		}

		if action.shouldRun(a) {
			go action.run(a)
		}

		a.requeueScheduledAction(action)
	}
}

func (a *App) popScheduledAction() scheduledAction {
	action, _ := a.scheduledActions.Pop()
	return action.(scheduledAction)
}

func (a *App) requeueScheduledAction(action scheduledAction) {
	action.updateNextRunTime(a)
	a.scheduledActions.Insert(action, float64(action.getNextRunTime().Unix()))
}
