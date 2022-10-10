package main

import (
	"fmt"
	"time"

	ga "github.com/saml-dev/gome-assistant"
)

func main() {
	app, err := ga.App("192.168.86.67:8123")
	defer app.Cleanup()
	if err != nil {
		fmt.Println(err)
		return
	}

	s := ga.ScheduleBuilder().Call(lightsOut).Daily().At(ga.HourMinute(22, 00)).Build()
	s2 := ga.ScheduleBuilder().Call(lightsOut).Every(time.Hour*4 + time.Minute*30).Offset(ga.HourMinute(1, 0)).Build()
	app.RegisterSchedule(s2)
	// err = app.Start()

	simpleListener := ga.EntityListenerBuilder().
		EntityId("light.lights").
		Call(cool).
		OnlyBetween(ga.HourMinute(22, 00), ga.HourMinute(07, 00))
	fmt.Println(simpleListener)

	fmt.Println(s, "\n", s2)
}

func lightsOut(service ga.Service) {
	// ga.TurnOff("light.all_lights")
}

func cool(service ga.Service, data ga.Data) {
	service.Light.TurnOn("light.entryway_lamp")
}
