package gomeassistant

import (
	"fmt"
	"time"

	"github.com/saml-dev/gome-assistant/internal"
)

type IntervalCallback func(*Service, *State)

type Interval struct {
	frequency   time.Duration
	callback    IntervalCallback
	startTime   TimeString
	endTime     TimeString
	nextRunTime time.Time

	exceptionDays   []time.Time
	exceptionRanges []timeRange
}

func (i Interval) Hash() string {
	return fmt.Sprint(i.startTime, i.endTime, i.frequency, i.callback, i.exceptionDays, i.exceptionRanges)
}

// Call
type intervalBuilder struct {
	interval Interval
}

// Every
type intervalBuilderCall struct {
	interval Interval
}

// Offset, ExceptionDay, ExceptionRange
type intervalBuilderEnd struct {
	interval Interval
}

func NewInterval() intervalBuilder {
	return intervalBuilder{
		Interval{
			frequency: 0,
			startTime: "00:00",
			endTime:   "00:00",
		},
	}
}

func (i Interval) String() string {
	return fmt.Sprintf("Interval{ call %q every %s%s%s }",
		internal.GetFunctionName(i.callback),
		i.frequency,
		formatStartOrEndString(i.startTime /* isStart = */, true),
		formatStartOrEndString(i.endTime /* isStart = */, false),
	)
}

func formatStartOrEndString(s TimeString, isStart bool) string {
	if s == "00:00" {
		return ""
	}
	if isStart {
		return fmt.Sprintf(" starting at %s", s)
	} else {
		return fmt.Sprintf(" ending at %s", s)
	}
}

func (sb intervalBuilder) Call(callback IntervalCallback) intervalBuilderCall {
	sb.interval.callback = callback
	return intervalBuilderCall(sb)
}

// Takes a DurationString ("2h", "5m", etc) to set the frequency of the interval.
func (sb intervalBuilderCall) Every(s DurationString) intervalBuilderEnd {
	d := internal.ParseDuration(string(s))
	sb.interval.frequency = d
	return intervalBuilderEnd(sb)
}

// Takes a TimeString ("HH:MM") when this interval will start running for the day.
func (sb intervalBuilderEnd) StartingAt(s TimeString) intervalBuilderEnd {
	sb.interval.startTime = s
	return sb
}

// Takes a TimeString ("HH:MM") when this interval will stop running for the day.
func (sb intervalBuilderEnd) EndingAt(s TimeString) intervalBuilderEnd {
	sb.interval.endTime = s
	return sb
}

func (sb intervalBuilderEnd) ExceptionDay(t time.Time) intervalBuilderEnd {
	sb.interval.exceptionDays = append(sb.interval.exceptionDays, t)
	return sb
}

func (sb intervalBuilderEnd) ExceptionRange(start, end time.Time) intervalBuilderEnd {
	sb.interval.exceptionRanges = append(sb.interval.exceptionRanges, timeRange{start, end})
	return sb
}

func (sb intervalBuilderEnd) Build() Interval {
	return sb.interval
}

// app.Start() functions
func runIntervals(a *App) {
	if a.intervals.Len() == 0 {
		return
	}

	for {
		i := popInterval(a)

		// run callback for all intervals before now in case they overlap
		for i.nextRunTime.Before(time.Now()) {
			i.maybeRunCallback(a)
			requeueInterval(a, i)

			i = popInterval(a)
		}

		time.Sleep(time.Until(i.nextRunTime))
		i.maybeRunCallback(a)
		requeueInterval(a, i)
	}
}

func (i Interval) maybeRunCallback(a *App) {
	if c := checkStartEndTime(i.startTime /* isStart = */, true); c.fail {
		return
	}
	if c := checkStartEndTime(i.endTime /* isStart = */, false); c.fail {
		return
	}
	if c := checkExceptionDays(i.exceptionDays); c.fail {
		return
	}
	if c := checkExceptionRanges(i.exceptionRanges); c.fail {
		return
	}
	go i.callback(a.service, a.state)
}

func popInterval(a *App) Interval {
	i, _ := a.intervals.Pop()
	return i.(Interval)
}

func requeueInterval(a *App, i Interval) {
	i.nextRunTime = i.nextRunTime.Add(i.frequency)

	a.intervals.Insert(i, float64(i.nextRunTime.Unix()))
}
