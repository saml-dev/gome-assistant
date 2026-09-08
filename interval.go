package gomeassistant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"saml.dev/gome-assistant/internal"
)

type IntervalCallback func(*Service, State)

// Interval is a `scheduledAction` that runs repeatedly at a specified
// frequency.
type Interval struct {
	frequency   time.Duration
	callback    IntervalCallback
	startTime   TimeString
	endTime     TimeString
	nextRunTime time.Time

	exceptionDates  []time.Time
	exceptionRanges []timeRange

	enabledEntities  []internal.EnabledDisabledInfo
	disabledEntities []internal.EnabledDisabledInfo

	key             string
	automationKey   string
	source          string
	metricsDisabled bool
}

func (i Interval) Hash() string {
	return fmt.Sprint(i.startTime, i.endTime, i.frequency, i.callback, i.exceptionDates, i.exceptionRanges)
}

// Call
type intervalBuilder struct {
	interval Interval
}

// Every
type intervalBuilderCall struct {
	interval Interval
}

// Offset, ExceptionDates, ExceptionRange
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

func (ib intervalBuilder) Call(callback IntervalCallback) intervalBuilderCall {
	ib.interval.callback = callback
	return intervalBuilderCall(ib)
}

// Takes a DurationString ("2h", "5m", etc) to set the frequency of the interval.
func (ib intervalBuilderCall) Every(s DurationString) intervalBuilderEnd {
	d := internal.ParseDuration(string(s))
	ib.interval.frequency = d
	return intervalBuilderEnd(ib)
}

// Takes a TimeString ("HH:MM") when this interval will start running for the day.
func (ib intervalBuilderEnd) StartingAt(s TimeString) intervalBuilderEnd {
	ib.interval.startTime = s
	return ib
}

// Takes a TimeString ("HH:MM") when this interval will stop running for the day.
func (ib intervalBuilderEnd) EndingAt(s TimeString) intervalBuilderEnd {
	ib.interval.endTime = s
	return ib
}

func (ib intervalBuilderEnd) ExceptionDates(t time.Time, tl ...time.Time) intervalBuilderEnd {
	ib.interval.exceptionDates = append(tl, t)
	return ib
}

func (ib intervalBuilderEnd) ExceptionRange(start, end time.Time) intervalBuilderEnd {
	ib.interval.exceptionRanges = append(ib.interval.exceptionRanges, timeRange{start, end})
	return ib
}

/*
Enable this interval only when the current state of {entityID} matches {state}.
If there is a network error while retrieving state, the interval runs if {runOnNetworkError} is true.
*/
func (ib intervalBuilderEnd) EnabledWhen(entityID, state string, runOnNetworkError bool) intervalBuilderEnd {
	if entityID == "" {
		panic(fmt.Sprintf("entityID is empty in EnabledWhen entityID='%s' state='%s'", entityID, state))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	ib.interval.enabledEntities = append(ib.interval.enabledEntities, i)
	return ib
}

/*
Disable this interval when the current state of {entityID} matches {state}.
If there is a network error while retrieving state, the interval runs if {runOnNetworkError} is true.
*/
func (ib intervalBuilderEnd) DisabledWhen(entityID, state string, runOnNetworkError bool) intervalBuilderEnd {
	if entityID == "" {
		panic(fmt.Sprintf("entityID is empty in EnabledWhen entityID='%s' state='%s'", entityID, state))
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	ib.interval.disabledEntities = append(ib.interval.disabledEntities, i)
	return ib
}

// Key sets the stable logical identity used by metrics and the web UI.
func (ib intervalBuilderEnd) Key(key string) intervalBuilderEnd {
	key = strings.TrimSpace(key)
	if key == "" {
		panic(ErrInvalidArgs)
	}
	ib.interval.key = key
	return ib
}

// DisableMetrics excludes this interval from metrics collection and the web
// UI catalog without disabling its execution.
func (ib intervalBuilderEnd) DisableMetrics() intervalBuilderEnd {
	ib.interval.metricsDisabled = true
	return ib
}

func (sb intervalBuilderEnd) Build() Interval {
	return sb.interval
}

// run invokes `i.maybeRunCallback()` based on its configured
// frequency.
func (i *Interval) run(ctx context.Context, app *App) {
	// Create a new, but stopped, timer for sleeping on:
	timer := time.NewTimer(1 * time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	defer stopAndDrainTimer(timer)

	for ctx.Err() == nil {
		if i.nextRunTime.After(time.Now()) {
			timer.Reset(time.Until(i.nextRunTime))
			select {
			case <-timer.C:
			case <-ctx.Done():
				return
			}
		}

		i.maybeRunCallback(app)
		i.updateNextRunTime(app)
	}
}

func (i *Interval) updateNextRunTime(app *App) {
	i.nextRunTime = i.nextRunTime.Add(i.frequency)
	app.updateCatalogNextDue(i.automationKey, i.nextRunTime)
}

func (i Interval) maybeRunCallback(app *App) {
	evaluation := app.newEvaluation(i.automationKey, intervalTrigger())
	if c := checkStartEndTime(i.startTime /* isStart = */, true); c.fail {
		evaluation.skip("before_start_time", nil, nil)
		return
	}
	if c := checkStartEndTime(i.endTime /* isStart = */, false); c.fail {
		evaluation.skip("after_end_time", nil, nil)
		return
	}
	if c := checkExceptionDates(i.exceptionDates); c.fail {
		evaluation.skip("exception_date", nil, nil)
		return
	}
	if c := checkExceptionRanges(i.exceptionRanges); c.fail {
		evaluation.skip("exception_range", nil, nil)
		return
	}
	if c := checkEnabledEntity(evaluation.scopedState(app.state), i.enabledEntities); c.fail {
		evaluation.skipCondition(c, "enabled_condition_mismatch", "enabled_condition_query_error")
		return
	}
	if c := checkDisabledEntity(evaluation.scopedState(app.state), i.disabledEntities); c.fail {
		evaluation.skipCondition(c, "disabled_condition_match", "disabled_condition_query_error")
		return
	}
	if app.sessionStopping() {
		evaluation.skip("session_stopping", nil, nil)
		return
	}
	if !evaluation.runCallback(app, func(service *Service, state State) { i.callback(service, state) }) {
		evaluation.skip("session_stopping", nil, nil)
	}
}
