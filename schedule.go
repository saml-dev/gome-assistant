package gomeassistant

import (
	"fmt"
	"time"

	"github.com/golang-module/carbon"
	"saml.dev/gome-assistant/internal"
)

type ScheduleCallback func(*Service, *State)

type DailySchedule struct {
	// 0-23
	hour int
	// 0-59
	minute int

	callback    ScheduleCallback
	nextRunTime time.Time

	isSunrise bool
	isSunset  bool
	sunOffset DurationString

	exceptionDates []time.Time
	allowlistDates []time.Time
}

func (s DailySchedule) Hash() string {
	return fmt.Sprint(s.hour, s.minute, s.callback)
}

type scheduleBuilder struct {
	schedule DailySchedule
}

type scheduleBuilderCall struct {
	schedule DailySchedule
}

type scheduleBuilderEnd struct {
	schedule DailySchedule
}

func NewDailySchedule() scheduleBuilder {
	return scheduleBuilder{
		DailySchedule{
			hour:      0,
			minute:    0,
			sunOffset: "0s",
		},
	}
}

func (s DailySchedule) String() string {
	return fmt.Sprintf("Schedule{ call %q daily at %s }",
		internal.GetFunctionName(s.callback),
		stringHourMinute(s.hour, s.minute),
	)
}

func stringHourMinute(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func (sb scheduleBuilder) Call(callback ScheduleCallback) scheduleBuilderCall {
	sb.schedule.callback = callback
	return scheduleBuilderCall(sb)
}

// At takes a string in 24hr format time like "15:30".
func (sb scheduleBuilderCall) At(s string) scheduleBuilderEnd {
	t := internal.ParseTime(s)
	sb.schedule.hour = t.Hour()
	sb.schedule.minute = t.Minute()
	return scheduleBuilderEnd(sb)
}

// Sunrise takes an optional duration string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (sb scheduleBuilderCall) Sunrise(offset ...DurationString) scheduleBuilderEnd {
	sb.schedule.isSunrise = true
	if len(offset) > 0 {
		sb.schedule.sunOffset = offset[0]
	}
	return scheduleBuilderEnd(sb)
}

// Sunset takes an optional duration string that is passed to time.ParseDuration.
// Examples include "-1.5h", "30m", etc. See https://pkg.go.dev/time#ParseDuration
// for full list.
func (sb scheduleBuilderCall) Sunset(offset ...DurationString) scheduleBuilderEnd {
	sb.schedule.isSunset = true
	if len(offset) > 0 {
		sb.schedule.sunOffset = offset[0]
	}
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderEnd) ExceptionDates(t time.Time, tl ...time.Time) scheduleBuilderEnd {
	sb.schedule.exceptionDates = append(tl, t)
	return sb
}

func (sb scheduleBuilderEnd) OnlyOnDates(t time.Time, tl ...time.Time) scheduleBuilderEnd {
	sb.schedule.allowlistDates = append(tl, t)
	return sb
}

func (sb scheduleBuilderEnd) Build() DailySchedule {
	return sb.schedule
}

// app.Start() functions
func runSchedules(a *App) {
	if a.schedules.Len() == 0 {
		return
	}

	for {
		sched := popSchedule(a)

		// run callback for all schedules before now in case they overlap
		for sched.nextRunTime.Before(time.Now()) {
			sched.maybeRunCallback(a)
			requeueSchedule(a, sched)

			sched = popSchedule(a)
		}

		fmt.Println("Next schedule:", sched.nextRunTime)
		time.Sleep(time.Until(sched.nextRunTime))
		sched.maybeRunCallback(a)
		requeueSchedule(a, sched)
	}
}

func (s DailySchedule) maybeRunCallback(a *App) {
	if c := checkExceptionDates(s.exceptionDates); c.fail {
		return
	}
	if c := checkAllowlistDates(s.allowlistDates); c.fail {
		return
	}
	go s.callback(a.service, a.state)
}

func popSchedule(a *App) DailySchedule {
	_sched, _ := a.schedules.Pop()
	return _sched.(DailySchedule)
}

func requeueSchedule(a *App, s DailySchedule) {
	if s.isSunrise || s.isSunset {
		var nextSunTime carbon.Carbon
		// "0s" is default value
		if s.sunOffset != "0s" {
			nextSunTime = getNextSunRiseOrSet(a, s.isSunrise, s.sunOffset)
		} else {
			nextSunTime = getNextSunRiseOrSet(a, s.isSunrise)
		}

		s.nextRunTime = nextSunTime.Carbon2Time()
	} else {
		s.nextRunTime = carbon.Time2Carbon(s.nextRunTime).AddDay().Carbon2Time()
	}

	a.schedules.Insert(s, float64(s.nextRunTime.Unix()))
}
