package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	ga "saml.dev/gome-assistant"
	"saml.dev/gome-assistant/app"
	gaapp "saml.dev/gome-assistant/app"
	"saml.dev/gome-assistant/websocket"
)

func main() {
	ctx := context.Background()
	app, err := gaapp.NewApp(
		ctx,
		gaapp.NewAppRequest{
			IpAddress:        "192.168.86.67", // Replace with your Home Assistant IP Address
			HAAuthToken:      os.Getenv("HA_AUTH_TOKEN"),
			HomeZoneEntityID: "zone.home",
		},
	)
	if err != nil {
		slog.Error("Error connecting to HASS:", err)
		os.Exit(1)
	}

	defer app.Close()

	pantryDoor := gaapp.
		NewEntityListener().
		EntityIDs("binary_sensor.pantry_door").
		Call(func(sensor gaapp.EntityData) {
			pantryLights(app, sensor)
		}).
		Build()

	_11pmSched := gaapp.
		NewDailySchedule().
		Call(func() {
			lightsOut(app)
		}).
		At("23:00").
		Build()

	_30minsBeforeSunrise := gaapp.
		NewDailySchedule().
		Call(func() {
			sunriseSched(app)
		}).
		Sunrise("-30m").
		Build()

	zwaveEventListener := gaapp.
		NewEventListener().
		EventTypes("zwave_js_value_notification").
		Call(onEvent).
		Build()

	app.RegisterEntityListeners(pantryDoor)
	app.RegisterSchedules(_11pmSched, _30minsBeforeSunrise)
	app.RegisterEventListeners(zwaveEventListener)

	app.Start(ctx)
}

func pantryLights(app *app.App, sensor gaapp.EntityData) {
	l := ga.EntityTarget("light.pantry")
	if sensor.ToState == "on" {
		app.Service.HomeAssistant.TurnOn(l, nil)
	} else {
		app.Service.HomeAssistant.TurnOff(l)
	}
}

func onEvent(ev websocket.Event) {
	// Since the structure of the event data changes depending on the
	// event type, you can Unmarshal the data into a Go type. If a
	// type for your event doesn't exist, you can write it yourself!
	// PR's welcome to the eventTypes.go file :)
	var data gaapp.ZWaveJSEventData
	json.Unmarshal(ev.RawData, &data)
	slog.Info("On event invoked", "data", data)
}

func lightsOut(app *app.App) {
	// always turn off outside lights
	app.Service.Light.TurnOff(ga.EntityTarget("light.outside_lights"))
	s, err := app.State.Get("binary_sensor.living_room_motion")
	if err != nil {
		slog.Warn("couldnt get living room motion state, doing nothing")
		return
	}

	// if no motion detected in living room for 30mins
	if s.State == "off" && time.Since(s.LastChanged).Minutes() > 30 {
		app.Service.Light.TurnOff(ga.EntityTarget("light.main_lights"))
	}
}

func sunriseSched(app *app.App) {
	app.Service.Light.TurnOn(ga.EntityTarget("light.living_room_lamps"), nil)
	app.Service.Light.TurnOff(ga.EntityTarget("light.christmas_lights"))
}
