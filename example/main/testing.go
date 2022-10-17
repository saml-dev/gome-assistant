package main

import (
	"fmt"
	"log"

	ga "github.com/saml-dev/gome-assistant"
)

func main() {
	app := ga.NewApp("192.168.86.67:8123")
	defer app.Cleanup()
	s := ga.ScheduleBuilder().Call(lightsOut).Daily().At(app.Sunset("1h")).Build()
	s2 := ga.ScheduleBuilder().Call(lightsOut).Every("2h").Offset("10m").Build()
	fmt.Println(s2)
	app.RegisterSchedule(s)
	app.RegisterSchedule(s2)
	simpleListener := ga.EntityListenerBuilder().
		EntityIds("group.office_ceiling_lights").
		Call(listenerCB).
		OnlyAfter("23:03").
		// Throttle("5s").
		Build()
	app.RegisterEntityListener(simpleListener)

	app.Start()

}

func lightsOut(service *ga.Service, state *ga.State) {
	// service.InputDatetime.Set("input_datetime.garage_last_triggered_ts", time.Now())
	// service.HomeAssistant.Toggle("group.living_room_lamps", map[string]any{"brightness_pct": 100})
	// service.Light.Toggle("light.entryway_lamp", map[string]any{"brightness_pct": 100})
	service.HomeAssistant.Toggle("light.entryway_lamp")
	log.Default().Println("running lightsOut")
	// service.HomeAssistant.Toggle("light.entryway_lamp")
	// log.Default().Println("A")
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
