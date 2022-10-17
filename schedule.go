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

// app.Start() functions
func RunSchedules(a *app) {
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

func popSchedule(a *app) schedule {
	_sched, _ := a.schedules.Pop()
	return _sched.(schedule)
}

func requeueSchedule(a *app, s schedule) {
	s.realStartTime = s.realStartTime.Add(s.frequency)
	a.schedules.Insert(s, float64(s.realStartTime.Unix()))
}
