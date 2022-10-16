package gomeassistant

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
)

type sunriseSunset struct {
	base        time.Duration
	addition    time.Duration
	subtraction time.Duration
}

func Sunrise() *sunriseSunset {
	return &sunriseSunset{
		base:        TimeOfDay(0, 10000),
		addition:    TimeOfDay(0, 0),
		subtraction: TimeOfDay(0, 0),
	}
}

func Sunset() *sunriseSunset {
	return &sunriseSunset{
		base:        TimeOfDay(0, 20000),
		addition:    TimeOfDay(0, 0),
		subtraction: TimeOfDay(0, 0),
	}
}

func (ss *sunriseSunset) Add(hm time.Duration) *sunriseSunset {
	ss.addition = hm
	return ss
}

func (ss *sunriseSunset) Subtract(hm time.Duration) *sunriseSunset {
	ss.subtraction = hm
	return ss
}

func (ss *sunriseSunset) Minutes() float64 {
	return ss.base.Minutes() + ss.addition.Minutes() - ss.subtraction.Minutes()
}

type timeOfDay interface {
	// Time represented as number of Minutes
	// after midnight. E.g. 02:00 would be 120.
	Minutes() float64
}

// TimeOfDay is a helper function to easily represent
// a time of day as a time.Duration since midnight.
func TimeOfDay(hour, minute int) time.Duration {
	return time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute)
}

// Duration is a wrapper for TimeOfDay that makes
// semantic sense when used with Every()
func Duration(hour, minute int) time.Duration {
	return TimeOfDay(hour, minute)
}

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
		offset is 4 character string representing hours and minutes
		in a 24-hr format.
		It is the base that your frequency will be added to.
		Defaults to "0000" (which is probably fine for most cases).

		Example: Run in the 3rd minute of every hour.
			Schedule{
				frequency: gomeassistant.Hourly // helper const for time.Hour
				offset: "0003"
			}
	*/
	offset time.Duration
	/*
		err will be set rather than returning an error to avoid checking err for nil on every schedule :)
		RegisterSchedule will exit if the error is set.
	*/
	err           error
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
			offset:    TimeOfDay(0, 0),
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

func (sb scheduleBuilderDaily) At(t timeOfDay) scheduleBuilderEnd {
	sb.schedule.offset = convertTimeOfDayToActualOffset(t)
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderCall) Every(duration time.Duration) scheduleBuilderCustom {
	sb.schedule.frequency = duration
	return scheduleBuilderCustom(sb)
}

func (sb scheduleBuilderCustom) Offset(t time.Duration) scheduleBuilderEnd {
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

func convertTimeOfDayToActualOffset(t timeOfDay) time.Duration {
	mins := t.Minutes()
	if mins > 15000 {
		// TODO: same as below but w/ sunset
		// don't forget to subtract 20000 here
		return TimeOfDay(0, 0)
	} else if mins > 5000 {
		// TODO: use httpClient to get state of sun.sun
		// to get next sunrise time
		// don't forget to subtract 10000 here to get +- from sunrise that user requested

		// retrieve next sunrise time

		// use carbon.Parse() to create time.Time of that time

		// return Time() of that many hours and minutes to set offset from midnight
	} else if mins >= 1440 {
		log.Fatalln("Offset (set via At() or Offset()) cannot be more than 1 day (23h59m)")
	}
	return TimeOfDay(0, int(mins))
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
