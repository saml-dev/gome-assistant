package app

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

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/internal/http"
	"saml.dev/gome-assistant/internal/priorityqueue"
	"saml.dev/gome-assistant/websocket"
)

// Returned by NewApp() if authentication fails
var ErrInvalidToken = websocket.ErrInvalidToken

var ErrInvalidArgs = errors.New("invalid arguments provided")

type App struct {
	// Wraps the ws connection with added mutex locking
	wsConn *websocket.Conn

	httpClient *http.HttpClient

	Service *Service
	State   State

	scheduledActions priorityqueue.PriorityQueue
	entityListeners  map[string][]*EntityListener
	eventListeners   map[string][]*EventListener

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
	// EntityID of the zone representing your home e.g. "zone.home".
	// Used to pull latitude/longitude from Home Assistant
	// to calculate sunset/sunrise times.
	HomeZoneEntityID string
}

// NewAppFromConfig establishes the websocket connection and returns
// an object you can use to register schedules and listeners, based on
// the URIs that it should connect to. `ctx` is used only to limit the
// time spent connecting; it cannot be used after that to cancel the
// app.
func NewAppFromConfig(ctx context.Context, config NewAppConfig) (*App, error) {
	if config.RESTBaseURI == "" || config.WebsocketURI == "" ||
		config.HAAuthToken == "" || config.HomeZoneEntityID == "" {
		slog.Error(
			"RESTBaseURI, WebsocketURI, HAAuthToken, and HomeZoneEntityID " +
				"are all required arguments in NewAppRequest",
		)
		return nil, ErrInvalidArgs
	}

	wsWriter, err := websocket.NewConnFromURI(ctx, config.WebsocketURI, config.HAAuthToken)
	if err != nil {
		return nil, err
	}

	httpClient := http.ClientFromUri(config.RESTBaseURI, config.HAAuthToken)

	state, err := newState(httpClient, config.HomeZoneEntityID)
	if err != nil {
		return nil, err
	}
	app := App{
		wsConn:           wsWriter,
		httpClient:       httpClient,
		State:            state,
		scheduledActions: priorityqueue.New(),
		entityListeners:  map[string][]*EntityListener{},
		eventListeners:   map[string][]*EventListener{},
		cancel:           func() {},
	}
	app.Service = newService(&app, httpClient)

	return &app, nil
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

// NewApp establishes the websocket connection and returns an object
// you can use to register schedules and listeners. `ctx` is used only
// to limit the time spent connecting; it cannot be used after that to
// cancel the app. If this function returns successfully, then
// `App.Close()` must eventually be called to release resources.
func NewApp(ctx context.Context, request NewAppRequest) (*App, error) {
	if request.IpAddress == "" || request.HAAuthToken == "" || request.HomeZoneEntityID == "" {
		slog.Error(
			"IpAddress, HAAuthToken, and HomeZoneEntityID " +
				"are all required arguments in NewAppRequest",
		)
		return nil, ErrInvalidArgs
	}
	port := request.Port
	if port == "" {
		port = "8123"
	}

	config := NewAppConfig{
		HAAuthToken:      request.HAAuthToken,
		HomeZoneEntityID: request.HomeZoneEntityID,
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
	initializeNextRunTime(app *App)
	shouldRun(app *App) bool
	run(app *App)
	updateNextRunTime(app *App)
	getNextRunTime() time.Time
}

func (app *App) RegisterScheduledAction(action scheduledAction) {
	action.initializeNextRunTime(app)
	app.scheduledActions.Insert(action, float64(action.getNextRunTime().Unix()))
}

func (app *App) RegisterSchedules(schedules ...*DailySchedule) {
	for _, s := range schedules {
		app.RegisterScheduledAction(s)
	}
}

func (app *App) RegisterIntervals(intervals ...*Interval) {
	for _, i := range intervals {
		app.RegisterScheduledAction(i)
	}
}

func (app *App) RegisterEntityListener(etl EntityListener) {
	if etl.delay != 0 && etl.toState == "" {
		slog.Error("EntityListener error: you have to use ToState() when using Duration()")
		panic(ErrInvalidArgs)
	}

	for _, entity := range etl.entityIDs {
		if elList, ok := app.entityListeners[entity]; ok {
			app.entityListeners[entity] = append(elList, &etl)
		} else {
			app.entityListeners[entity] = []*EntityListener{&etl}
		}
	}
}

func (app *App) RegisterEntityListeners(etls ...EntityListener) {
	for _, etl := range etls {
		app.RegisterEntityListener(etl)
	}
}

func (app *App) RegisterEventListener(evl EventListener) {
	for _, eventType := range evl.eventTypes {
		elList, ok := app.eventListeners[eventType]
		if !ok {
			// FIXME: keep track of subscriptions so that they can
			// be unsubscribed from.
			_, err := app.WatchEvents(
				eventType,
				func(msg websocket.Message) {
					go app.callEventListeners(msg)
				},
			)
			if err != nil {
				// FIXME: better error handling
				panic(err)
			}
		}
		app.eventListeners[eventType] = append(elList, &evl)
	}
}

func (app *App) RegisterEventListeners(evls ...EventListener) {
	for _, evl := range evls {
		app.RegisterEventListener(evl)
	}
}

func getSunriseSunset(
	s State, sunrise bool, dateToUse carbon.Carbon, offset ...DurationString,
) carbon.Carbon {
	date := dateToUse.Carbon2Time()
	rise, set := sunriseLib.SunriseSunset(
		s.Latitude(), s.Longitude(), date.Year(), date.Month(), date.Day(),
	)
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
			parsingErr := fmt.Errorf(
				"could not parse offset passed to %s: \"%s\": %w",
				printString, offset[0], err,
			)
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
	sunriseOrSunset := getSunriseSunset(app.State, sunrise, carbon.Now(), offset...)
	if sunriseOrSunset.Lt(carbon.Now()) {
		// if we're past today's sunset or sunrise (accounting for offset) then get tomorrows
		// as that's the next time the schedule will run
		sunriseOrSunset = getSunriseSunset(app.State, sunrise, carbon.Tomorrow(), offset...)
	}
	return sunriseOrSunset
}

// Start the app. When `ctx` expires, the app closes the connection
// and returns.
func (app *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	app.cancel = cancel
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	slog.Info("Starting", "scheduled actions", app.scheduledActions.Len())
	slog.Info("Starting", "entity listeners", len(app.entityListeners))
	slog.Info("Starting", "event listeners", len(app.eventListeners))

	eg.Go(func() error {
		app.runScheduledActions(ctx)
		return nil
	})

	// subscribe to state_changed events
	stateChangedSubscription, err := app.WatchStateChangedEvents(
		func(msg websocket.Message) {
			go app.callEntityListeners(msg)
		},
	)
	if err != nil {
		return fmt.Errorf("subscribing to 'state_changed' events: %w", err)
	}

	defer app.unwatchEvents(stateChangedSubscription)

	// entity listeners runOnStartup
	for eid, etls := range app.entityListeners {
		eid := eid
		for _, etl := range etls {
			etl := etl
			// ensure each ETL only runs once, even if
			// it listens to multiple entities
			if etl.runOnStartup && !etl.runOnStartupCompleted {
				entityState, err := app.State.Get(eid)
				if err != nil {
					slog.Warn(
						"Failed to get entity state \"", eid,
						"\" during startup, skipping RunOnStartup",
					)
				}

				etl.runOnStartupCompleted = true
				eg.Go(func() error {
					etl.callback(EntityData{
						TriggerEntityID: eid,
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
	eg.Go(func() error {
		app.wsConn.Start()
		cancel()
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		app.Close()
		return nil
	})

	eg.Wait()

	return nil
}

// Close closes the connection and releases any resources. It may be
// called more than once; only the first call does anything.
func (app *App) Close() {
	app.closeOnce.Do(func() {
		app.close()
	})
}

// close closes the connection and releases resources. It must be
// called exactly once.
func (app *App) close() {
	app.cancel()
	app.wsConn.Close()
}

func (app *App) GetService() *Service {
	return app.Service
}

func (app *App) GetState() State {
	return app.State
}

func (app *App) runScheduledActions(ctx context.Context) {
	if app.scheduledActions.Len() == 0 {
		return
	}

	// Create a new, but stopped, timer:
	timer := time.NewTimer(1 * time.Hour)
	if !timer.Stop() {
		<-timer.C
	}

	for {
		action := app.popScheduledAction()
		if action.getNextRunTime().After(time.Now()) {
			timer.Reset(time.Until(action.getNextRunTime()))

			select {
			case <-timer.C:
			case <-ctx.Done():
				return
			}
		}

		if action.shouldRun(app) {
			go action.run(app)
		}

		app.requeueScheduledAction(action)
	}
}

func (app *App) popScheduledAction() scheduledAction {
	action, _ := app.scheduledActions.Pop()
	return action.(scheduledAction)
}

func (app *App) requeueScheduledAction(action scheduledAction) {
	action.updateNextRunTime(app)
	app.scheduledActions.Insert(action, float64(action.getNextRunTime().Unix()))
}

type SubEvent struct {
	websocket.BaseMessage
	EventType string `json:"event_type"`
}

// WatchEvents subscribes to events of the given type, invoking
// `subscriber` when any such events are received. Calls to
// `subscriber` are synchronous with respect to any other received
// messages, but asynchronous with respect to writes.
func (app *App) WatchEvents(
	eventType string, subscriber websocket.Subscriber,
) (websocket.Subscription, error) {
	// Make sure we're listening before events might start arriving:
	e := SubEvent{
		BaseMessage: websocket.BaseMessage{
			Type: "subscribe_events",
		},
		EventType: eventType,
	}
	var subscription websocket.Subscription
	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		subscription = lc.Subscribe(subscriber)
		e.ID = subscription.ID()
		if err := lc.SendMessage(e); err != nil {
			lc.Unsubscribe(subscription)
			return fmt.Errorf("error writing to websocket: %w", err)
		}
		return nil
	})
	if err != nil {
		return websocket.Subscription{}, err
	}
	// m, _ := ReadMessage(conn, ctx)
	// log.Default().Println(string(m))
	return subscription, nil
}

func (app *App) WatchStateChangedEvents(
	subscriber websocket.Subscriber,
) (websocket.Subscription, error) {
	return app.WatchEvents("state_changed", subscriber)
}

type UnsubEvent struct {
	websocket.BaseMessage
	Subscription int64 `json:"subscription"`
}

// unwatchEvents unsubscribes to events with the given `subscriptionID`. This does
// not remove the subscriber.
func (app *App) unwatchEvents(subscription websocket.Subscription) error {
	e := UnsubEvent{
		BaseMessage: websocket.BaseMessage{
			Type: "unsubscribe_events",
		},
		Subscription: subscription.ID(),
	}

	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		lc.Unsubscribe(subscription)

		e.ID = lc.NextID()
		return lc.SendMessage(e)
	})
	if err != nil {
		return fmt.Errorf("unsubscribing from ID %d: %w", subscription.ID(), err)
	}

	// m, _ := ReadMessage(conn, ctx)
	// log.Default().Println(string(m))
	return nil
}

// Call invokes an RPC service corresponding to `req` via websockets
// and waits for and returns a single `result`. `msg` must be
// serializable to JSON. It shouldn't have its ID filled in yet; that
// will be done within this method. The response is not analyzed at
// all, even to check for errors.
func (app *App) Call(
	ctx context.Context, req websocket.Request,
) (websocket.Message, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseCh := make(chan websocket.Message, 1)

	var subscription websocket.Subscription

	// Receive a single message, sent it to `responseCh`, then
	// unsubscribe:
	subscriber := func(msg websocket.Message) {
		defer close(responseCh)
		responseCh <- msg
		_ = app.wsConn.Send(func(lc websocket.LockedConn) error {
			lc.Unsubscribe(subscription)
			return nil
		})
	}

	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		subscription = lc.Subscribe(subscriber)
		req.SetID(subscription.ID())
		if err := lc.SendMessage(req); err != nil {
			lc.Unsubscribe(subscription)
			return fmt.Errorf("error writing to websocket: %w", err)
		}
		return nil
	})

	if err != nil {
		return websocket.Message{}, err
	}

	select {
	case response := <-responseCh:
		return response, nil
	case <-ctx.Done():
		return websocket.Message{}, ctx.Err()
	}
}

type CallServiceRequest struct {
	websocket.BaseMessage
	Domain  string `json:"domain"`
	Service string `json:"service"`

	// ServiceData must be serializable to a JSON object.
	ServiceData any `json:"service_data,omitempty"`

	Target ga.Target `json:"target,omitempty"`
}

// CallService invokes a service using a `call_service` message, then
// waits for and returns the response.
//
// FIXME: can the response be parsed into a result-style message?
func (app *App) CallService(
	ctx context.Context, domain string, service string, serviceData any, target ga.Target,
) (websocket.Message, error) {
	req := CallServiceRequest{
		BaseMessage: websocket.BaseMessage{
			Type: "call_service",
		},
		Domain:      domain,
		Service:     service,
		ServiceData: serviceData,
		Target:      target,
	}

	return app.Call(ctx, &req)
}

// Subscribe subscribes to some events via `req`, waits for a single
// response, and then leaves `subscriber` subscribed to the events. If
// this method returns without an error, `subscriber` must eventually
// be unsubscribed. `ctx` covers the subscription and the wait for the
// first answer, but not the forwarding of subsequent events or
// unsubscribing.
//
// FIXME: should this subscriber and subscription be specialized to
// event messages?
//
// FIXME: should the result be examined? If the subscription request
// failed, then we could fail more generally instead of leaving the
// cleanup to the caller.
func (app *App) Subscribe(
	ctx context.Context, req websocket.Request, subscriber websocket.Subscriber,
) (websocket.Message, websocket.Subscription, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// The result of the attempt to subscribe (i.e., the first
	// message) will be sent to this channel.
	resultReceived := false
	resultCh := make(chan websocket.Message, 1)

	var subscription websocket.Subscription

	// Receive a single message, sent it to `responseCh`, then
	// unsubscribe:
	dualSubscriber := func(msg websocket.Message) {
		if !resultReceived {
			// This is the first message. We send it to the channel so
			// that it can be returned from the outer function.
			defer close(resultCh)
			resultCh <- msg
			resultReceived = true
			return
		}

		// The result has already been processed. Subsequent events
		// get forwarded to `subscriber`:
		subscriber(msg)
	}

	err := app.wsConn.Send(func(lc websocket.LockedConn) error {
		subscription = lc.Subscribe(dualSubscriber)
		req.SetID(subscription.ID())
		if err := lc.SendMessage(req); err != nil {
			lc.Unsubscribe(subscription)
			return fmt.Errorf("error writing to websocket: %w", err)
		}
		return nil
	})

	if err != nil {
		return websocket.Message{}, websocket.Subscription{}, err
	}

	select {
	case response := <-resultCh:
		return response, subscription, nil
	case <-ctx.Done():
		return websocket.Message{}, websocket.Subscription{}, ctx.Err()
	}
}
