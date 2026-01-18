package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saml.dev/gome-assistant/cmd/example/entities" // Optional import generated entities
	"saml.dev/gome-assistant/message"

	ga "saml.dev/gome-assistant"
)

//go:generate go run saml.dev/gome-assistant/cmd/generate

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app, err := ga.NewApp(
		ctx,
		ga.NewAppRequest{
			URL:              "http://192.168.86.67:8123", // Replace with your Home Assistant URL
			HAAuthToken:      os.Getenv("HA_AUTH_TOKEN"),
			HomeZoneEntityID: "zone.home",
		},
	)
	if err != nil {
		slog.Error("Error connecting to HASS:", "error", err)
		os.Exit(1)
	}

	defer app.Cleanup()

	pantryDoor := ga.
		NewEntityListener().
		EntityIDs(entities.BinarySensor.PantryDoor). // Use generated entity constant
		Call(func(service *ga.Service, state ga.State, sensor message.StateChangedData) {
			pantryLights(ctx, service, state, sensor)
		}).
		Build()

	_11pmSched := ga.
		NewDailySchedule().
		Call(func(service *ga.Service, state ga.State) {
			lightsOut(ctx, service, state)
		}).
		At("23:00").
		Build()

	_30minsBeforeSunrise := ga.
		NewDailySchedule().
		Call(func(service *ga.Service, state ga.State) {
			sunriseSched(ctx, service, state)
		}).
		Sunrise("-30m").
		Build()

	zwaveEventListener := ga.
		NewEventListener().
		EventTypes("zwave_js_value_notification").
		Call(onEvent).
		Build()

	app.RegisterEntityListeners(pantryDoor)
	app.RegisterSchedules(_11pmSched, _30minsBeforeSunrise)
	app.RegisterEventListeners(zwaveEventListener)

	app.Start()
}

func pantryLights(
	ctx context.Context, service *ga.Service, state ga.State, sensor message.StateChangedData,
) {
	l := "light.pantry"
	// l := entities.Light.Pantry // Or use generated entity constant
	if sensor.NewState.State == "on" {
		if _, err := service.HomeAssistant.TurnOn(ctx, l); err != nil {
			slog.Warn("couldn't turn on pantry light")
		}
	} else {
		if _, err := service.HomeAssistant.TurnOff(ctx, l); err != nil {
			slog.Warn("couldn't turn off pantry light")
		}
	}
}

func onEvent(service *ga.Service, state ga.State, msg message.Message) {
	// Since the structure of the event changes depending
	// on the event type, you can Unmarshal the raw json
	// into a Go type. If a type for your event doesn't
	// exist, you can write it yourself! PR's welcome to
	// the eventTypes.go file :)
	ev := ga.EventZWaveJSValueNotification{}
	json.Unmarshal(msg.Raw, &ev)
	slog.Info("On event invoked", "event", ev)
}

func lightsOut(ctx context.Context, service *ga.Service, state ga.State) {
	// always turn off outside lights
	if _, err := service.Light.TurnOff(ctx, entities.Light.OutsideLights); err != nil {
		slog.Warn("couldn't turn off living room light, doing nothing")
		return
	}
	s, err := state.Get(entities.BinarySensor.LivingRoomMotion)
	if err != nil {
		slog.Warn("couldn't get living room motion state, doing nothing")
		return
	}

	// if no motion detected in living room for 30mins
	if s.State == "off" && time.Since(s.LastChanged).Minutes() > 30 {
		if _, err := service.Light.TurnOff(ctx, entities.Light.MainLights); err != nil {
			slog.Warn("couldn't turn off living light")
			return
		}
	}
}

func sunriseSched(ctx context.Context, service *ga.Service, state ga.State) {
	if _, err := service.Light.TurnOn(ctx, entities.Light.LivingRoomLamps); err != nil {
		slog.Warn("couldn't turn on living light")
	}

	if _, err := service.Light.TurnOff(ctx, entities.Light.ChristmasLights); err != nil {
		slog.Warn("couldn't turn off Christmas lights")
	}
}
