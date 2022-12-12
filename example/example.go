package main

import (
	"encoding/json"
	"log"
	"time"

	ga "saml.dev/gome-assistant"
)

func main() {
	app := ga.NewApp("0.0.0.0:8123") // Replace with your Home Assistant IP Address
	defer app.Cleanup()

	pantryDoor := ga.
		NewEntityListener().
		EntityIds("binary_sensor.pantry_door").
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

func pantryLights(service *ga.Service, state *ga.State, sensor ga.EntityData) {
	l := "light.pantry"
	if sensor.ToState == "on" {
		service.HomeAssistant.TurnOn(l)
	} else {
		service.HomeAssistant.TurnOff(l)
	}
}

func onEvent(service *ga.Service, state *ga.State, data ga.EventData) {
	// Since the structure of the event changes depending
	// on the event type, you can Unmarshal the raw json
	// into a Go type. If a type for your event doesn't
	// exist, you can write it yourself! PR's welcome to
	// the eventTypes.go file :)
	ev := ga.EventZWaveJSValueNotification{}
	json.Unmarshal(data.RawEventJSON, &ev)
	log.Default().Println(ev)
}

func lightsOut(service *ga.Service, state *ga.State) {
	// always turn off outside lights
	service.Light.TurnOff("light.outside_lights")
	s, err := state.Get("binary_sensor.living_room_motion")
	if err != nil {
		log.Default().Println("couldnt get living room motion state, doing nothing")
		return
	}

	// if no motion detected in living room for 30mins
	if s.State == "off" && time.Since(s.LastChanged).Minutes() > 30 {
		service.Light.TurnOff("light.main_lights")
	}
}

func sunriseSched(service *ga.Service, state *ga.State) {
	service.Light.TurnOn("light.living_room_lamps")
	service.Light.TurnOff("light.christmas_lights")
}
