package gomeassistant

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
)

type sunriseSunset struct {
	base        timeOfDay
	addition    timeOfDay
	subtraction timeOfDay
}

var Sunrise *sunriseSunset = &sunriseSunset{
	base:        TimeOfDay(0, 10000),
	addition:    TimeOfDay(0, 0),
	subtraction: TimeOfDay(0, 0),
}

var Sunset *sunriseSunset = &sunriseSunset{
	base:        TimeOfDay(0, 20000),
	addition:    TimeOfDay(0, 0),
	subtraction: TimeOfDay(0, 0),
}

func (ss *sunriseSunset) Add(hm timeOfDay) *sunriseSunset {
	ss.addition = hm
	return ss
}

func (ss *sunriseSunset) Subtract(hm timeOfDay) *sunriseSunset {
	ss.subtraction = hm
	return ss
}

func (ss *sunriseSunset) minutes() int {
	return ss.base.minute +
		(ss.addition.hour*60 + ss.addition.minute) -
		(ss.subtraction.hour*60 + ss.subtraction.minute)
}

// HourMinute is used to express a time of day
// but it shouldn't be used directly. Use
// HourMinute(), Sunset(), or Sunrise() to
// create one. Add() and Subtract() can be
// called on Sunset and Sunrise to offset
// from that time.
type timeOfDay struct {
	hour   int
	minute int
}

type timeOfDayInterface interface {
	// Time represented as number of minutes
	// after midnight. E.g. 02:00 would be 120.
	minutes() int
}

func (hm timeOfDay) minutes() int {
	return hm.hour*60 + hm.minute
}

func TimeOfDay(Hour, Minute int) timeOfDay {
	return timeOfDay{Hour, Minute}
}

func (hm timeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", hm.hour, hm.minute)
}

type scheduleCallback func(Service, State)

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
	offset timeOfDay
	/*
		This will be set rather than returning an error to avoid checking err for nil on every schedule :)
		RegisterSchedule will exit if the error is set.
	*/
	err           error
	realStartTime time.Time
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
			offset:    timeOfDay{0, 0},
		},
	}
}

func (s schedule) String() string {
	return fmt.Sprintf("Run %q %s %s",
		getFunctionName(s.callback),
		frequencyToString(s.frequency),
		s.offset,
	)
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

func (sb scheduleBuilderDaily) At(t timeOfDayInterface) scheduleBuilderEnd {
	sb.schedule.offset = convertTimeOfDayToActualOffset(t)
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderCall) Every(duration time.Duration) scheduleBuilderCustom {
	sb.schedule.frequency = duration
	return scheduleBuilderCustom(sb)
}

func (sb scheduleBuilderCustom) Offset(t timeOfDay) scheduleBuilderEnd {
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

func convertTimeOfDayToActualOffset(t timeOfDayInterface) timeOfDay {
	if t.minutes() > 15000 {
		// TODO: same as below but w/ sunset
		return TimeOfDay(0, 0)
	} else if t.minutes() > 5000 {
		// TODO: use httpClient to get state of sun.sun
		// to get next sunrise time

		// retrieve next sunrise time

		// use carbon.Parse() to create time.Time of that time

		// return Time() of that many hours and minutes to set offset from midnight
	} else if t.minutes() >= 1440 {
		log.Fatalln("Offset (set via At() or Offset()) cannot be more than 1 day (23h59m)")
	}
	return TimeOfDay(0, t.minutes())
}
