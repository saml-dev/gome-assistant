package main

import (
	"log"

	ga "github.com/saml-dev/gome-assistant"
)

func main() {
	app := ga.NewApp("192.168.86.67:8123")
	defer app.Cleanup()
	pantryDoor := ga.
		EntityListenerBuilder().
		EntityIds("binary_sensor.pantry_door").
		Call(pantryLights).
		Build()
	app.RegisterEntityListener(pantryDoor)

	app.Start()

}

func pantryLights(service *ga.Service, data ga.EntityData) {
	// service.InputDatetime.Set("input_datetime.garage_last_triggered_ts", time.Now())
	// service.HomeAssistant.Toggle("group.living_room_lamps", map[string]any{"brightness_pct": 100})
	// service.Light.Toggle("light.entryway_lamp", map[string]any{"brightness_pct": 100})
	if data.ToState == "on" {
		service.HomeAssistant.TurnOn("switch.pantry_light_2")
	} else {
		service.HomeAssistant.TurnOff("switch.pantry_light_2")
	}
}

func cool(service *ga.Service, state *ga.State) {
	// service.Light.TurnOn("light.entryway_lamp")
	// log.Default().Println("B")
}

func c(service *ga.Service, state *ga.State) {
	// log.Default().Println("C")
}

func listenerCB(service *ga.Service, data ga.EntityData) {
	log.Default().Println("hi katie")
}

// TODO: randomly placed, add .Throttle to Listener
