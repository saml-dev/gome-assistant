package gomeassistant

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/saml-dev/gome-assistant/internal"
)

type scheduleCallback func(*Service, *State)

type schedule struct {
	/*
		frequency is a time.Duration representing how often you want to run your function.

		Some examples:
			time.Second * 5 // runs every 5 seconds at 00:00:00, 00:00:05, etc.
			time.Hour * 12 // runs at offset, +12 hours, +24 hours, etc.
			gomeassistant.Daily // runs at offset, +24 hours, +48 hours, etc. Daily is a const helper for time.Hour * 24
			// Helpers include Daily, Hourly, Minutely
	*/
	frequency time.Duration
	callback  scheduleCallback
	/*
		offset is the base that your frequency will be added to.
		Defaults to 0 (which is probably fine for most cases).

		Example: Run in the 3rd minute of every hour.
			ScheduleBuilder().Call(myFunc).Every("1h").Offset("3m")
	*/
	offset        time.Duration
	realStartTime time.Time
}

func (s schedule) Hash() string {
	return fmt.Sprint(s.offset, s.frequency, s.callback)
}

type scheduleBuilder struct {
	schedule schedule
}

type scheduleBuilderCall struct {
	schedule schedule
}

type scheduleBuilderDaily struct {
	schedule schedule
}

type scheduleBuilderCustom struct {
	schedule schedule
}

type scheduleBuilderEnd struct {
	schedule schedule
}

func ScheduleBuilder() scheduleBuilder {
	return scheduleBuilder{
		schedule{
			frequency: 0,
			offset:    0,
		},
	}
}

func (s schedule) String() string {
	return fmt.Sprintf("Schedule{ call %q %s %s }",
		getFunctionName(s.callback),
		frequencyToString(s.frequency),
		offsetToString(s),
	)
}

func offsetToString(s schedule) string {
	if s.frequency.Hours() == 24 {
		return fmt.Sprintf("%02d:%02d", int(s.offset.Hours()), int(s.offset.Minutes())%60)
	}
	return s.offset.String()
}

func frequencyToString(d time.Duration) string {
	if d.Hours() == 24 {
		return "daily at"
	}
	return "every " + d.String() + " with offset"
}

func (sb scheduleBuilder) Call(callback scheduleCallback) scheduleBuilderCall {
	sb.schedule.callback = callback
	return scheduleBuilderCall(sb)
}

func (sb scheduleBuilderCall) Daily() scheduleBuilderDaily {
	sb.schedule.frequency = time.Hour * 24
	return scheduleBuilderDaily(sb)
}

// At takes a string 24hr format time like "15:30".
func (sb scheduleBuilderDaily) At(s string) scheduleBuilderEnd {
	t := internal.ParseTime(s)
	sb.schedule.offset = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderDaily) Sunrise(s string) scheduleBuilderEnd {
	// TODO: this function should calculate sunrise time,
	// set isSunrise on schedule,
	// return sb.At() with the string value caluclated via sunrise and offset time

	// for calculating next sunset with -30m offset, just do today's sunset + 24 hours. It'll only ever be like a minute off

	// since we have isSunrise flag and it's daily schedule, could use offset field on schedule
	// to hold optional offset here
	// NOTE: this doesn't work, At() IS setting the offset with a duration of 24h. Prolly need to set sunRiseSetOffset or something on schedule

	// NOTE: below line can't work because needs connString and auth token for http client.
	// Maybe just pass app to NewSchedule or put NewSchedule/NewEntList/NewEvList
	// ON app itself? Then it has access to app through the builders?
	// sunriseString := getSunriseSunset()
	// d, err :=
}

func (sb scheduleBuilderCall) Every(s TimeString) scheduleBuilderCustom {
	d, err := time.ParseDuration(string(s))
	if err != nil {
		log.Fatalf("couldn't parse string duration passed to Every(): \"%s\" see https://pkg.go.dev/time#ParseDuration for valid time units", s)
	}
	sb.schedule.frequency = d
	return scheduleBuilderCustom(sb)
}

func (sb scheduleBuilderCustom) Offset(s TimeString) scheduleBuilderEnd {
	t, err := time.ParseDuration(string(s))
	if err != nil {
		log.Fatalf("Couldn't parse string duration passed to Offset(): \"%s\" see https://pkg.go.dev/time#ParseDuration for valid time units", s)
	}
	sb.schedule.offset = t
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderCustom) Build() schedule {
	return sb.schedule
}

func (sb scheduleBuilderEnd) Build() schedule {
	return sb.schedule
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type SunsetSchedule struct {
	offset TimeString
}

// app.Start() functions
func runSchedules(a *app) {
	if a.schedules.Len() == 0 {
		return
	}

	for {
		sched := popSchedule(a)
		// log.Default().Println(sched.realStartTime)

		// run callback for all schedules before now in case they overlap
		for sched.realStartTime.Before(time.Now()) {
			go sched.callback(a.service, a.state)
			requeueSchedule(a, sched)

			sched = popSchedule(a)
		}

		time.Sleep(time.Until(sched.realStartTime))
		go sched.callback(a.service, a.state)
		requeueSchedule(a, sched)
	}
}

func popSchedule(a *app) ScheduleInterface {
	_sched, _ := a.schedules.Pop()
	return _sched.(ScheduleInterface)
}

func requeueSchedule(a *app, s ScheduleInterface) {
	nextTime := s.GetNext()
	a.schedules.Insert(s, float64(nextTime.Unix()))
}
