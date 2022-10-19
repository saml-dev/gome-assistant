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
	zwaveEventListener := ga.
		EventListenerBuilder().
		EventType("zwave_js_value_notification").
		Call(onEvent).
		Build()
	app.RegisterEntityListener(pantryDoor)
	app.RegisterSchedule(ga.ScheduleBuilder().Call(cool).Every("5s").Build())
	app.RegisterEventListener(zwaveEventListener)

	app.Start()

}

func pantryLights(service *ga.Service, data ga.EntityData) {
	l := "group.kitchen_ceiling_lights"
	// service.HomeAssistant.Toggle("group.living_room_lamps", map[string]any{"brightness_pct": 100})
	// service.Light.Toggle("light.entryway_lamp", map[string]any{"brightness_pct": 100})
	if data.ToState == "on" {
		service.HomeAssistant.TurnOn(l)
	} else {
		service.HomeAssistant.TurnOff(l)
	}
}

func onEvent(service *ga.Service, data ga.EventData) {
	service.HomeAssistant.Toggle("light.el_gato_key_lights")
}

func cool(service *ga.Service, state *ga.State) {
	// service.InputDatetime.Set("input_datetime.garage_last_triggered_ts", time.Now())
	// service.Light.TurnOn("light.entryway_lamp")
	// log.Default().Println("B")
}

func c(service *ga.Service, state *ga.State) {
	// log.Default().Println("C")
}

func listenerCB(service *ga.Service, data ga.EntityData) {
	log.Default().Println("hi")
}
