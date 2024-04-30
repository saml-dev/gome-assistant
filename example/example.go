package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	ga "saml.dev/gome-assistant"
	gaapp "saml.dev/gome-assistant/app"
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
		Call(pantryLights).
		Build()

	_11pmSched := gaapp.
		NewDailySchedule().
		Call(lightsOut).
		At("23:00").
		Build()

	_30minsBeforeSunrise := gaapp.
		NewDailySchedule().
		Call(sunriseSched).
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

func pantryLights(service *gaapp.Service, state gaapp.State, sensor gaapp.EntityData) {
	l := ga.EntityTarget("light.pantry")
	if sensor.ToState == "on" {
		service.HomeAssistant.TurnOn(l, nil)
	} else {
		service.HomeAssistant.TurnOff(l)
	}
}

func onEvent(service *gaapp.Service, state gaapp.State, data gaapp.EventData) {
	// Since the structure of the event changes depending
	// on the event type, you can Unmarshal the raw json
	// into a Go type. If a type for your event doesn't
	// exist, you can write it yourself! PR's welcome to
	// the eventTypes.go file :)
	ev := gaapp.EventZWaveJSValueNotification{}
	json.Unmarshal(data.RawEventJSON, &ev)
	slog.Info("On event invoked", "event", ev)
}

func lightsOut(service *gaapp.Service, state gaapp.State) {
	// always turn off outside lights
	service.Light.TurnOff(ga.EntityTarget("light.outside_lights"))
	s, err := state.Get("binary_sensor.living_room_motion")
	if err != nil {
		slog.Warn("couldnt get living room motion state, doing nothing")
		return
	}

	// if no motion detected in living room for 30mins
	if s.State == "off" && time.Since(s.LastChanged).Minutes() > 30 {
		service.Light.TurnOff(ga.EntityTarget("light.main_lights"))
	}
}

func sunriseSched(service *gaapp.Service, state gaapp.State) {
	service.Light.TurnOn(ga.EntityTarget("light.living_room_lamps"), nil)
	service.Light.TurnOff(ga.EntityTarget("light.christmas_lights"))
}
