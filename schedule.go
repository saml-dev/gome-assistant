package gomeassistant

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type hourMinute struct {
	Hour   int
	Minute int
}

func HourMinute(Hour, Minute int) hourMinute {
	return hourMinute{Hour, Minute}
}

func (hm hourMinute) String() string {
	return fmt.Sprintf("%02d:%02d", hm.Hour, hm.Minute)
}

type scheduleCallback func(Service)

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
	offset hourMinute
	/*
		This will be set rather than returning an error to avoid checking err for nil on every schedule :)
		RegisterSchedule will panic if the error is set.
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
	return scheduleBuilder{schedule{}}
}

func (s schedule) String() string {
	return fmt.Sprintf("Run %q %s %s",
		getFunctionName(s.callback),
		frequencyToString(s.frequency),
		s.offset,
	)
}

func frequencyToString(d time.Duration) string {
	fmt.Println(d.Hours(), d.Minutes(), d.Seconds())
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

func (sb scheduleBuilderDaily) At(t hourMinute) scheduleBuilderEnd {
	sb.schedule.offset = t
	return scheduleBuilderEnd(sb)
}

func (sb scheduleBuilderCall) Every(d time.Duration) scheduleBuilderCustom {
	sb.schedule.frequency = d
	return scheduleBuilderCustom(sb)
}

func (sb scheduleBuilderCustom) Offset(o hourMinute) scheduleBuilderEnd {
	sb.schedule.offset = o
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
