package gomeassistant

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	"saml.dev/gome-assistant/internal"
)

type ScheduleCallback func(*Service, State)

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

	enabledEntities  []internal.EnabledDisabledInfo
	disabledEntities []internal.EnabledDisabledInfo
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

/*
Enable this schedule only when the current state of {entityId} matches {state}.
If there is a network error while retrieving state, the schedule runs if {runOnNetworkError} is true.
*/
func (sb scheduleBuilderEnd) EnabledWhen(entityId, state string, runOnNetworkError bool) scheduleBuilderEnd {
	if entityId == "" {
		panic(fmt.Sprintf("entityId is empty in EnabledWhen entityId='%s' state='%s'", entityId, state))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityId,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	sb.schedule.enabledEntities = append(sb.schedule.enabledEntities, i)
	return sb
}

/*
Disable this schedule when the current state of {entityId} matches {state}.
If there is a network error while retrieving state, the schedule runs if {runOnNetworkError} is true.
*/
func (sb scheduleBuilderEnd) DisabledWhen(entityId, state string, runOnNetworkError bool) scheduleBuilderEnd {
	if entityId == "" {
		panic(fmt.Sprintf("entityId is empty in EnabledWhen entityId='%s' state='%s'", entityId, state))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityId,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	sb.schedule.disabledEntities = append(sb.schedule.disabledEntities, i)
	return sb
}

func (sb scheduleBuilderEnd) Build() DailySchedule {
	return sb.schedule
}

// app.Start() functions
func (a *App) runSchedules() {
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, sched := range a.schedules {
		wg.Add(1)
		go func(sched DailySchedule) {
			defer wg.Done()
			sched.run(a)
		}(sched)
	}
}

// run invokes `s.maybeRunCallback()` based on its configured
// schedule.
func (s DailySchedule) run(a *App) {
	for {
		if s.nextRunTime.After(time.Now()) {
			slog.Info("Next schedule", "start_time", s.nextRunTime)
			time.Sleep(time.Until(s.nextRunTime))
		}

		s.maybeRunCallback(a)
		s.updateNextRunTime(a)
	}
}

func (s DailySchedule) maybeRunCallback(a *App) {
	if c := checkExceptionDates(s.exceptionDates); c.fail {
		return
	}
	if c := checkAllowlistDates(s.allowlistDates); c.fail {
		return
	}
	if c := checkEnabledEntity(a.state, s.enabledEntities); c.fail {
		return
	}
	if c := checkDisabledEntity(a.state, s.disabledEntities); c.fail {
		return
	}
	go s.callback(a.service, a.state)
}

// updateNextRunTime updates `s.nextRunTime` to the next time that `s`
// should run.
func (s DailySchedule) updateNextRunTime(a *App) {
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
}
