package main

import (
	"log"
	"time"

	ga "github.com/saml-dev/gome-assistant"
)

func main() {
	app := ga.NewApp("192.168.86.67:8123")
	defer app.Cleanup()
	s := ga.ScheduleBuilder().Call(lightsOut).Every(time.Second * 5).Build()
	s2 := ga.ScheduleBuilder().Call(cool).Every(time.Millisecond * 500).Build()
	s3 := ga.ScheduleBuilder().Call(c).Every(time.Minute * 1).Build()
	app.RegisterSchedule(s)
	app.RegisterSchedule(s2)
	app.RegisterSchedule(s3)
	app.Start()

	simpleListener := ga.EntityListenerBuilder().
		EntityIds("light.lights").
		Call(listenerCB).
		OnlyBetween(ga.TimeOfDay(22, 00), ga.TimeOfDay(07, 00))
	log.Println(simpleListener)

	log.Println(s)
	log.Println(s2)
}

func lightsOut(service *ga.Service, state *ga.State) {
	// service.InputDatetime.Set("input_datetime.garage_last_triggered_ts", time.Now())
	// service.HomeAssistant.Toggle("group.living_room_lamps", map[string]any{"brightness_pct": 100})
	// service.Light.Toggle("light.entryway_lamp", map[string]any{"brightness_pct": 100})
	// service.Switch.Toggle("switch.espurna_sunroom_lamp")
	log.Default().Println("A")
}

func cool(service *ga.Service, state *ga.State) {
	// service.Light.TurnOn("light.entryway_lamp")
	log.Default().Println("B")
}

func c(service *ga.Service, state *ga.State) {
	log.Default().Println("C")
}

func listenerCB(service *ga.Service, data *ga.Data) {}
