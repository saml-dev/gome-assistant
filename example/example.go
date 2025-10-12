package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"example/entities" // Optional import generated entities

	ga "saml.dev/gome-assistant"
)

//go:generate go run saml.dev/gome-assistant/cmd/generate

func main() {
	ctx := context.TODO()

	app, err := ga.NewApp(
		ctx,
		ga.NewAppRequest{
			URL:              "http://192.168.86.67:8123", // Replace with your Home Assistant URL
			HAAuthToken:      os.Getenv("HA_AUTH_TOKEN"),
			HomeZoneEntityId: "zone.home",
		},
	)
	if err != nil {
		slog.Error("Error connecting to HASS:", "error", err)
		os.Exit(1)
	}

	defer app.Cleanup()

	pantryDoor := ga.
		NewEntityListener().
		EntityIds(entities.BinarySensor.PantryDoor). // Use generated entity constant
		Call(pantryLights).
		Build()

	_11pmSched := ga.
		NewDailySchedule().
		Call(lightsOut).
		At("23:00").
		Build()

	_30minsBeforeSunrise := ga.
		NewDailySchedule().
		Call(sunriseSched).
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

func pantryLights(service *ga.Service, state ga.State, sensor ga.EntityData) {
	l := "light.pantry"
	// l := entities.Light.Pantry // Or use generated entity constant
	if sensor.ToState == "on" {
		service.HomeAssistant.TurnOn(l)
	} else {
		service.HomeAssistant.TurnOff(l)
	}
}

func onEvent(service *ga.Service, state ga.State, data ga.EventData) {
	// Since the structure of the event changes depending
	// on the event type, you can Unmarshal the raw json
	// into a Go type. If a type for your event doesn't
	// exist, you can write it yourself! PR's welcome to
	// the eventTypes.go file :)
	ev := ga.EventZWaveJSValueNotification{}
	json.Unmarshal(data.RawEventJSON, &ev)
	slog.Info("On event invoked", "event", ev)
}

func lightsOut(service *ga.Service, state ga.State) {
	// always turn off outside lights
	service.Light.TurnOff(entities.Light.OutsideLights)
	s, err := state.Get(entities.BinarySensor.LivingRoomMotion)
	if err != nil {
		slog.Warn("couldnt get living room motion state, doing nothing")
		return
	}

	// if no motion detected in living room for 30mins
	if s.State == "off" && time.Since(s.LastChanged).Minutes() > 30 {
		service.Light.TurnOff(entities.Light.MainLights)
	}
}

func sunriseSched(service *ga.Service, state ga.State) {
	service.Light.TurnOn(entities.Light.LivingRoomLamps)
	service.Light.TurnOff(entities.Light.ChristmasLights)
}
