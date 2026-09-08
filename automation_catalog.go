package gomeassistant

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"

	"saml.dev/gome-assistant/internal"
)

type automationDefinition struct {
	Key            string
	Type           string
	Source         string
	Callback       string
	TriggerSummary string
	NextDue        time.Time
	LastEvaluation time.Time
	LastOutcome    string
	EvaluationID   string
	Configuration  map[string]any
}

// automationCatalogSnapshot returns the frozen definitions with the small set
// of runtime fields that the eventual web layer needs. Definitions themselves
// are never persisted or revised, and callers cannot mutate App-owned slices
// or maps through the returned value.
func (app *App) automationCatalogSnapshot() []automationDefinition {
	app.configurationMu.Lock()
	defer app.configurationMu.Unlock()

	definitions := make([]automationDefinition, len(app.catalog))
	for index, definition := range app.catalog {
		definitions[index] = definition
		definitions[index].Configuration = cloneConfiguration(definition.Configuration)
	}
	return definitions
}

func cloneConfiguration(configuration map[string]any) map[string]any {
	if configuration == nil {
		return nil
	}
	copy := make(map[string]any, len(configuration))
	for key, value := range configuration {
		switch value := value.(type) {
		case []string:
			copy[key] = append([]string(nil), value...)
		case []map[string]string:
			items := make([]map[string]string, len(value))
			for index, item := range value {
				items[index] = mapsClone(item)
			}
			copy[key] = items
		case []map[string]any:
			items := make([]map[string]any, len(value))
			for index, item := range value {
				items[index] = mapsClone(item)
			}
			copy[key] = items
		default:
			copy[key] = value
		}
	}
	return copy
}

func mapsClone[T any](source map[string]T) map[string]T {
	copy := make(map[string]T, len(source))
	for key, value := range source {
		copy[key] = value
	}
	return copy
}

func (app *App) updateCatalogStatus(key, evaluationID, outcome string, occurredAt time.Time) {
	if key == "" {
		return
	}
	app.configurationMu.Lock()
	defer app.configurationMu.Unlock()
	for index := range app.catalog {
		if app.catalog[index].Key != key {
			continue
		}
		app.catalog[index].LastEvaluation = occurredAt
		app.catalog[index].LastOutcome = outcome
		app.catalog[index].EvaluationID = evaluationID
		return
	}
}

func (app *App) updateCatalogNextDue(key string, nextDue time.Time) {
	if key == "" {
		return
	}
	app.configurationMu.Lock()
	defer app.configurationMu.Unlock()
	for index := range app.catalog {
		if app.catalog[index].Key == key {
			app.catalog[index].NextDue = nextDue
			return
		}
	}
}

func callbackLocation(callback any) (source, name string, err error) {
	value := reflect.ValueOf(callback)
	if !value.IsValid() || value.Kind() != reflect.Func || value.IsNil() {
		return "", "", fmt.Errorf("callback source metadata is unavailable; supply an explicit Key")
	}
	fn := runtime.FuncForPC(value.Pointer())
	if fn == nil {
		return "", "", fmt.Errorf("callback source metadata is unavailable; supply an explicit Key")
	}
	file, line := fn.FileLine(value.Pointer())
	if file == "" || line <= 0 {
		return "", fn.Name(), fmt.Errorf("callback source metadata is unavailable; supply an explicit Key")
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line), fn.Name(), nil
}

func resolveAutomationIdentity(explicit string, callback any) (key, source, name string, err error) {
	source, name, sourceErr := callbackLocation(callback)
	if explicit != "" {
		return explicit, source, name, nil
	}
	if sourceErr != nil {
		return "", source, name, sourceErr
	}
	return source, source, name, nil
}

func (app *App) buildAutomationCatalog() ([]automationDefinition, error) {
	definitions := make([]automationDefinition, 0, len(app.scheduledActions)+len(app.entityListeners)+len(app.eventListeners))

	for _, action := range app.scheduledActions {
		switch automation := action.(type) {
		case *DailySchedule:
			if automation.metricsDisabled {
				automation.automationKey, automation.source = "", ""
				continue
			}
			key, source, callback, err := resolveAutomationIdentity(automation.key, automation.callback)
			if err != nil {
				return nil, fmt.Errorf("daily automation: %w", err)
			}
			automation.automationKey, automation.source = key, source
			definitions = append(definitions, dailyDefinition(automation, key, source, callback))
		case *Interval:
			if automation.metricsDisabled {
				automation.automationKey, automation.source = "", ""
				continue
			}
			key, source, callback, err := resolveAutomationIdentity(automation.key, automation.callback)
			if err != nil {
				return nil, fmt.Errorf("interval automation: %w", err)
			}
			automation.automationKey, automation.source = key, source
			definitions = append(definitions, intervalDefinition(automation, key, source, callback))
		}
	}

	seenEntities := make(map[*EntityListener]struct{})
	for _, listeners := range app.entityListeners {
		for _, automation := range listeners {
			if _, seen := seenEntities[automation]; seen {
				continue
			}
			seenEntities[automation] = struct{}{}
			if automation.metricsDisabled {
				automation.automationKey, automation.source = "", ""
				continue
			}
			key, source, callback, err := resolveAutomationIdentity(automation.key, automation.callback)
			if err != nil {
				return nil, fmt.Errorf("entity automation: %w", err)
			}
			automation.automationKey, automation.source = key, source
			definitions = append(definitions, entityDefinition(automation, key, source, callback))
		}
	}

	seenEvents := make(map[*EventListener]struct{})
	for _, listeners := range app.eventListeners {
		for _, automation := range listeners {
			if _, seen := seenEvents[automation]; seen {
				continue
			}
			seenEvents[automation] = struct{}{}
			if automation.metricsDisabled {
				automation.automationKey, automation.source = "", ""
				continue
			}
			key, source, callback, err := resolveAutomationIdentity(automation.key, automation.callback)
			if err != nil {
				return nil, fmt.Errorf("event automation: %w", err)
			}
			automation.automationKey, automation.source = key, source
			definitions = append(definitions, eventDefinition(automation, key, source, callback))
		}
	}

	owners := make(map[string]automationDefinition, len(definitions))
	for _, definition := range definitions {
		if previous, exists := owners[definition.Key]; exists {
			return nil, fmt.Errorf("duplicate automation key %q used by %s (%s) and %s (%s); callback source locations are used as keys by default; use the .Key() builder method to give each automation a unique key", definition.Key, previous.Type, previous.Source, definition.Type, definition.Source)
		}
		owners[definition.Key] = definition
	}
	sort.Slice(definitions, func(i, j int) bool {
		if definitions[i].Type == definitions[j].Type {
			return definitions[i].Key < definitions[j].Key
		}
		return definitions[i].Type < definitions[j].Type
	})
	return definitions, nil
}

func dailyDefinition(a *DailySchedule, key, source, callback string) automationDefinition {
	trigger := stringHourMinute(a.hour, a.minute)
	if a.isSunrise {
		trigger = "sunrise"
	} else if a.isSunset {
		trigger = "sunset"
	}
	if (a.isSunrise || a.isSunset) && a.sunOffset != "0s" {
		trigger += " " + string(a.sunOffset)
	}
	return automationDefinition{
		Key: key, Type: "daily", Source: source, Callback: callback, TriggerSummary: trigger, NextDue: a.nextRunTime,
		Configuration: map[string]any{
			"time": stringHourMinute(a.hour, a.minute), "sunrise": a.isSunrise, "sunset": a.isSunset,
			"solar_offset": string(a.sunOffset), "exception_dates": formatDates(a.exceptionDates),
			"allowed_dates": formatDates(a.allowlistDates), "enabled_conditions": formatConditions(a.enabledEntities),
			"disabled_conditions": formatConditions(a.disabledEntities),
		},
	}
}

func intervalDefinition(a *Interval, key, source, callback string) automationDefinition {
	return automationDefinition{
		Key: key, Type: "interval", Source: source, Callback: callback,
		TriggerSummary: fmt.Sprintf("every %s%s%s", a.frequency, formatStartOrEndString(a.startTime, true), formatStartOrEndString(a.endTime, false)),
		NextDue:        a.nextRunTime,
		Configuration: map[string]any{
			"interval": a.frequency.String(), "start_time": string(a.startTime), "end_time": string(a.endTime),
			"exception_dates": formatDates(a.exceptionDates), "exception_ranges": formatRanges(a.exceptionRanges),
			"enabled_conditions": formatConditions(a.enabledEntities), "disabled_conditions": formatConditions(a.disabledEntities),
		},
	}
}

func entityDefinition(a *EntityListener, key, source, callback string) automationDefinition {
	return automationDefinition{
		Key: key, Type: "entity", Source: source, Callback: callback, TriggerSummary: strings.Join(a.entityIDs, ", "),
		Configuration: map[string]any{
			"entity_ids": append([]string(nil), a.entityIDs...), "start_time": a.betweenStart, "end_time": a.betweenEnd,
			"from_state": a.fromState, "to_state": a.toState, "duration": a.delay.String(), "throttle": a.throttle.String(),
			"exception_dates": formatDates(a.exceptionDates), "exception_ranges": formatRanges(a.exceptionRanges),
			"run_on_startup": a.runOnStartup, "ignore_attribute_changes": a.ignoreAttributeChanges,
			"enabled_conditions": formatConditions(a.enabledEntities), "disabled_conditions": formatConditions(a.disabledEntities),
		},
	}
}

func eventDefinition(a *EventListener, key, source, callback string) automationDefinition {
	return automationDefinition{
		Key: key, Type: "event", Source: source, Callback: callback, TriggerSummary: strings.Join(a.eventTypes, ", "),
		Configuration: map[string]any{
			"event_types": append([]string(nil), a.eventTypes...), "start_time": a.betweenStart, "end_time": a.betweenEnd,
			"throttle": a.throttle.String(), "exception_dates": formatDates(a.exceptionDates),
			"exception_ranges": formatRanges(a.exceptionRanges), "enabled_conditions": formatConditions(a.enabledEntities),
			"disabled_conditions": formatConditions(a.disabledEntities), "trigger_when": a.triggerWhen != nil,
		},
	}
}

func formatDates(dates []time.Time) []string {
	values := make([]string, len(dates))
	for i, date := range dates {
		values[i] = date.Format(time.RFC3339Nano)
	}
	return values
}

func formatRanges(ranges []timeRange) []map[string]string {
	values := make([]map[string]string, len(ranges))
	for i, value := range ranges {
		values[i] = map[string]string{"start": value.start.Format(time.RFC3339Nano), "end": value.end.Format(time.RFC3339Nano)}
	}
	return values
}

func formatConditions(conditions []internal.EnabledDisabledInfo) []map[string]any {
	values := make([]map[string]any, len(conditions))
	for i, condition := range conditions {
		values[i] = map[string]any{
			"entity_id": condition.Entity, "expected_state": condition.State, "run_on_network_error": condition.RunOnError,
		}
	}
	return values
}
