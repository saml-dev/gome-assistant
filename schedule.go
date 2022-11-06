package gomeassistant

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/saml-dev/gome-assistant/internal"
)

type ScheduleCallback func(*Service, *State)

type Schedule struct {
	frequency     time.Duration
	callback      ScheduleCallback
	offset        time.Duration
	realStartTime time.Time

	isSunrise bool
	isSunset  bool
	sunOffset DurationString

	exceptionDays   []time.Time
	exceptionRanges []timeRange
}

func (s Schedule) Hash() string {
	return fmt.Sprint(s.offset, s.frequency, s.callback)
}

type scheduleBuilder struct {
	schedule Schedule
}

type scheduleBuilderCall struct {
	schedule Schedule
}

type scheduleBuilderDaily struct {
	schedule Schedule
}

type scheduleBuilderCustom struct {
	schedule Schedule
}

type scheduleBuilderEnd struct {
	schedule Schedule
}

func NewSchedule() scheduleBuilder {
	return scheduleBuilder{
		Schedule{
			frequency: 0,
			offset:    0,
		},
	}
}

func (s Schedule) String() string {
	return fmt.Sprintf("Schedule{ call %q %s %s }",
		getFunctionName(s.callback),
		frequencyToString(s.frequency),
		offsetToString(s),
	)
}

func offsetToString(s Schedule) string {
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

func (sb scheduleBuilder) Call(callback ScheduleCallback) scheduleBuilderCall {
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

// Sunrise takes an app pointer and an optional duration string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (sb scheduleBuilderDaily) Sunrise(a *App, offset ...DurationString) scheduleBuilderEnd {
	sb.schedule.realStartTime = getSunriseSunsetFromApp(a, true, offset...).Carbon2Time()
	sb.schedule.isSunrise = true
	return scheduleBuilderEnd(sb)
}

// Sunset takes an app pointer and an optional duration string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (sb scheduleBuilderDaily) Sunset(a *App, offset ...DurationString) scheduleBuilderEnd {
	sb.schedule.realStartTime = getSunriseSunsetFromApp(a, false, offset...).Carbon2Time()
	sb.schedule.isSunset = true
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderCall) Every(s DurationString) scheduleBuilderCustom {
	d := internal.ParseDuration(string(s))
	sb.schedule.frequency = d
	return scheduleBuilderCustom(sb)
}

func (sb scheduleBuilderCustom) Offset(s DurationString) scheduleBuilderEnd {
	d := internal.ParseDuration(string(s))
	sb.schedule.offset = d
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderCustom) ExceptionDay(t time.Time) scheduleBuilderCustom {
	sb.schedule.exceptionDays = append(sb.schedule.exceptionDays, t)
	return sb
}

func (sb scheduleBuilderCustom) ExceptionRange(start, end time.Time) scheduleBuilderCustom {
	sb.schedule.exceptionRanges = append(sb.schedule.exceptionRanges, timeRange{start, end})
	return sb
}

func (sb scheduleBuilderCustom) Build() Schedule {
	return sb.schedule
}

func (sb scheduleBuilderEnd) ExceptionDay(t time.Time) scheduleBuilderEnd {
	sb.schedule.exceptionDays = append(sb.schedule.exceptionDays, t)
	return sb
}

func (sb scheduleBuilderEnd) ExceptionRange(start, end time.Time) scheduleBuilderEnd {
	sb.schedule.exceptionRanges = append(sb.schedule.exceptionRanges, timeRange{start, end})
	return sb
}

func (sb scheduleBuilderEnd) Build() Schedule {
	return sb.schedule
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// app.Start() functions
func runSchedules(a *App) {
	if a.schedules.Len() == 0 {
		return
	}

	for {
		sched := popSchedule(a)

		// run callback for all schedules before now in case they overlap
		for sched.realStartTime.Before(time.Now()) {
			maybeRunCallback(a, sched)
			requeueSchedule(a, sched)

			sched = popSchedule(a)
		}

		fmt.Println("Next schedule:", sched.realStartTime)
		time.Sleep(time.Until(sched.realStartTime))
		maybeRunCallback(a, sched)
		requeueSchedule(a, sched)
	}
}

func maybeRunCallback(a *App, s Schedule) {
	if c := checkExceptionDays(s.exceptionDays); c.fail {
		return
	}
	if c := checkExceptionRanges(s.exceptionRanges); c.fail {
		return
	}
	go s.callback(a.service, a.state)
}

func popSchedule(a *App) Schedule {
	_sched, _ := a.schedules.Pop()
	return _sched.(Schedule)
}

func requeueSchedule(a *App, s Schedule) {
	if s.isSunrise || s.isSunset {
		nextSunTime := getSunriseSunsetFromApp(a, s.isSunrise, s.sunOffset)

		// this is true when there is a negative offset, so schedule runs before sunset/sunrise and
		// HA still shows today's sunset as next sunset. Just add 24h as a default handler
		// since we can't get tomorrow's sunset from HA at this point.
		if nextSunTime.IsToday() {
			nextSunTime = nextSunTime.AddHours(24)
		}

		s.realStartTime = nextSunTime.Carbon2Time()
	} else {
		s.realStartTime = s.realStartTime.Add(s.frequency)
	}

	a.schedules.Insert(s, float64(s.realStartTime.Unix()))
}
