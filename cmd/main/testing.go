package main

import (
	"log"

	ga "github.com/saml-dev/gome-assistant"
)

func main() {
	app := ga.App("192.168.86.67:8123")
	defer app.Cleanup()
	s := ga.ScheduleBuilder().Call(lightsOut).Daily().At(ga.TimeOfDay(23, 00)).Build()
	s2 := ga.ScheduleBuilder().Call(lightsOut).Every(ga.Duration(04, 30)).Offset(ga.TimeOfDay(1, 0)).Build()
	app.RegisterSchedule(s2)
	app.Start()

	simpleListener := ga.EntityListenerBuilder().
		EntityId("light.lights").
		Call(cool).
		OnlyBetween(ga.TimeOfDay(22, 00), ga.TimeOfDay(07, 00))
	log.Println(simpleListener)

	log.Println(s)
	log.Println(s2)
}

func lightsOut(service ga.Service, state ga.State) {
	// ga.TurnOff("light.all_lights")
}

func cool(service ga.Service, data ga.Data) {
	service.Light.TurnOn("light.entryway_lamp")
}
