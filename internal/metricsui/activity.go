package metricsui

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"saml.dev/gome-assistant/internal/metrics"
)

type activityModel struct {
	Query    activityQuery
	Snapshot Snapshot
	Catalog  []Automation
	Rows     []activityRow
	HasOlder bool
}

type activityRow struct {
	ID         int64
	Evaluation *evaluationRow
	System     *systemRow
}

type evaluationRow struct {
	ID             string
	FirstID        int64
	LastID         int64
	SessionID      string
	AutomationKey  string
	AutomationType string
	Source         string
	Callback       string
	OccurredAt     time.Time
	LastOccurredAt time.Time
	Trigger        string
	Outcome        string
	Duration       time.Duration
	Warnings       int
	Incomplete     bool
	Trace          []traceItem
}

type systemRow struct {
	Event metrics.StoredEvent
	Text  string
}

type traceItem struct {
	Event metrics.StoredEvent
	Text  string
	Error string
}

func (ui *ui) activity(ctx context.Context, query activityQuery) (activityModel, error) {
	model := activityModel{Query: query, Snapshot: ui.options.Snapshot(), Catalog: ui.catalog()}
	if ui.options.Source == nil || ui.options.Source.Reader() == nil {
		return model, fmt.Errorf("metrics reader is unavailable")
	}
	reader := ui.options.Source.Reader()
	definitions := make(map[string]Automation, len(model.Catalog))
	for _, automation := range model.Catalog {
		definitions[automation.Key] = automation
	}

	rows := make([]activityRow, 0, activityScanLimit*2)
	if query.Scope != "system" {
		summaries, err := reader.Evaluations(ctx, metrics.EvaluationFilter{
			BeforeID:      query.BeforeID,
			AutomationKey: query.AutomationKey,
			Limit:         activityScanLimit,
		})
		if err != nil {
			return model, err
		}
		for _, summary := range summaries {
			trace, err := reader.EvaluationTrace(ctx, summary.EvaluationID)
			if err != nil {
				return model, err
			}
			row := summarizeEvaluation(summary, trace, definitions, model.Snapshot, ui.options.Now())
			if !matchesEvaluation(row, query) {
				continue
			}
			rows = append(rows, activityRow{ID: row.LastID, Evaluation: &row})
		}
	}
	if query.Scope != "automations" && query.Outcome == "" && query.AutomationType == "" && query.AutomationKey == "" {
		events, err := reader.Events(ctx, metrics.EventFilter{
			BeforeID: query.BeforeID,
			Scope:    metrics.ScopeSystem,
			Limit:    activityScanLimit,
		})
		if err != nil {
			return model, err
		}
		for _, event := range events {
			row := systemRow{Event: event, Text: describeSystemEvent(event)}
			if query.Search != "" && !containsFold(row.Text+" "+string(event.Type)+" "+string(event.PayloadJSON), query.Search) {
				continue
			}
			rows = append(rows, activityRow{ID: event.ID, System: &row})
		}
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].ID > rows[j].ID })
	// A page of exactly 100 logical rows is complete. Only offer pagination
	// when the merged evaluation/system result has a real 101st row.
	model.HasOlder = len(rows) > activityPageSize
	if len(rows) > activityPageSize {
		rows = rows[:activityPageSize]
	}
	model.Rows = rows
	return model, nil
}

func summarizeEvaluation(summary metrics.EvaluationSummary, trace []metrics.StoredEvent, definitions map[string]Automation, snapshot Snapshot, now time.Time) evaluationRow {
	definition := definitions[summary.AutomationKey]
	row := evaluationRow{
		ID:             summary.EvaluationID,
		FirstID:        summary.FirstID,
		LastID:         summary.LastID,
		SessionID:      summary.SessionID,
		AutomationKey:  summary.AutomationKey,
		AutomationType: definition.Type,
		Source:         definition.Source,
		Callback:       definition.Callback,
		Trigger:        definition.TriggerSummary,
		OccurredAt:     summary.FirstOccurred,
		LastOccurredAt: summary.LastOccurred,
		Duration:       summary.LastOccurred.Sub(summary.FirstOccurred),
		Trace:          make([]traceItem, 0, len(trace)),
	}
	if row.Source == "" {
		row.Source = "Current definition unavailable"
	}
	for _, event := range trace {
		item := traceItem{Event: event, Text: describeTraceEvent(event)}
		if errorText := payloadError(event.PayloadJSON); errorText != "" && isOperationFinished(event.Type) {
			row.Warnings++
			item.Error = errorText
		}
		row.Trace = append(row.Trace, item)
		switch event.Type {
		case metrics.EventAutomationTriggered:
			kind := triggerType(event.PayloadJSON)
			if row.AutomationType == "" {
				row.AutomationType = kind
			}
			// Current schedule definitions contain the useful concrete time/solar
			// summary; their persisted trigger intentionally contains only its kind.
			if row.Trigger == "" || (kind != "daily" && kind != "interval") {
				row.Trigger = describeTrigger(event.PayloadJSON)
			}
		case metrics.EventAutomationSkipped:
			row.Outcome = "skipped"
		case metrics.EventAutomationWaiting:
			row.Outcome = "waiting"
		case metrics.EventAutomationWaitCanceled:
			row.Outcome = "canceled"
		case metrics.EventAutomationCallbackStarted:
			row.Outcome = "running"
		case metrics.EventAutomationCallbackFinished:
			row.Outcome = "completed"
		}
	}
	if row.Trigger == "" {
		row.Trigger = "trigger details unavailable"
	}
	if row.AutomationType == "" {
		row.AutomationType = "unknown"
	}
	if row.Outcome != "completed" && row.Outcome != "skipped" && row.Outcome != "canceled" {
		// An unterminated trace (including one currently marked waiting) can be
		// live only when it belongs to the currently active session. A finished
		// (even most recent) session has no live callback to infer, so its missing
		// terminal fact is incomplete.
		if !snapshot.SessionActive || row.SessionID == "" || row.SessionID != snapshot.SessionID {
			row.Outcome = "incomplete"
			row.Incomplete = true
		} else if row.Outcome == "" {
			row.Outcome = "running"
		}
	}
	if row.Outcome == "running" || row.Outcome == "waiting" {
		row.Duration = now.Sub(row.OccurredAt)
		if row.Duration < 0 {
			row.Duration = 0
		}
	}
	return row
}

func matchesEvaluation(row evaluationRow, query activityQuery) bool {
	if query.AutomationKey != "" && row.AutomationKey != query.AutomationKey {
		return false
	}
	if query.Outcome != "" && row.Outcome != query.Outcome {
		return false
	}
	if query.AutomationType != "" && row.AutomationType != query.AutomationType {
		return false
	}
	return query.Search == "" || containsFold(strings.Join([]string{
		row.AutomationKey, row.Source, row.Callback, row.Trigger, row.Outcome, row.AutomationType,
	}, " "), query.Search)
}

func (ui *ui) catalog() []Automation {
	items := append([]Automation(nil), ui.options.Automations()...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type < items[j].Type
		}
		return compareCaseInsensitive(items[i].Key, items[j].Key) < 0
	})
	return items
}

func compareCaseInsensitive(left, right string) int {
	if order := strings.Compare(strings.ToLower(left), strings.ToLower(right)); order != 0 {
		return order
	}
	return strings.Compare(left, right)
}

func containsFold(value, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(strings.TrimSpace(needle)))
}

func triggerType(payload json.RawMessage) string {
	var body struct {
		Trigger struct {
			Kind string `json:"kind"`
		} `json:"trigger"`
	}
	_ = json.Unmarshal(payload, &body)
	if body.Trigger.Kind == "entity_startup" {
		return "entity"
	}
	return body.Trigger.Kind
}

func describeTrigger(payload json.RawMessage) string {
	var body struct {
		Trigger struct {
			Kind   string `json:"kind"`
			Entity *struct {
				ID   string `json:"trigger_entity_id"`
				From string `json:"from_state"`
				To   string `json:"to_state"`
			} `json:"entity"`
			Event *struct {
				Type string `json:"type"`
			} `json:"event"`
		} `json:"trigger"`
	}
	if err := json.Unmarshal(payload, &body); err != nil {
		return "trigger details unavailable"
	}
	switch body.Trigger.Kind {
	case "daily":
		return "daily schedule"
	case "interval":
		return "interval tick"
	case "entity", "entity_startup":
		if body.Trigger.Entity == nil {
			return "entity change"
		}
		if body.Trigger.Kind == "entity_startup" {
			return "startup state: " + body.Trigger.Entity.ID
		}
		return strings.TrimSpace(body.Trigger.Entity.ID + " " + body.Trigger.Entity.From + " → " + body.Trigger.Entity.To)
	case "event":
		if body.Trigger.Event != nil && body.Trigger.Event.Type != "" {
			return "event: " + body.Trigger.Event.Type
		}
		return "Home Assistant event"
	default:
		return "trigger details unavailable"
	}
}

func describeTraceEvent(event metrics.StoredEvent) string {
	switch event.Type {
	case metrics.EventAutomationTriggered:
		return "Triggered: " + describeTrigger(event.PayloadJSON)
	case metrics.EventAutomationSkipped:
		return describeSkip(event.PayloadJSON)
	case metrics.EventAutomationWaiting:
		return describeWait(event.PayloadJSON)
	case metrics.EventAutomationWaitCanceled:
		return "Duration wait canceled: " + payloadString(event.PayloadJSON, "reason")
	case metrics.EventAutomationCallbackStarted:
		return "Callback started"
	case metrics.EventAutomationCallbackFinished:
		return describeCallbackFinished(event.PayloadJSON)
	case metrics.EventStateReadStarted:
		return describeStateRead(event.PayloadJSON, "started")
	case metrics.EventStateReadFinished:
		return describeStateRead(event.PayloadJSON, "finished")
	case metrics.EventServiceCallStarted:
		return describeServiceOperation(event.PayloadJSON, "started")
	case metrics.EventServiceCallFinished:
		return describeServiceOperation(event.PayloadJSON, "finished")
	default:
		return strings.ReplaceAll(string(event.Type), "_", " ")
	}
}

func describeCallbackFinished(payload json.RawMessage) string {
	var body struct {
		DurationNS *int64 `json:"duration_ns"`
	}
	if json.Unmarshal(payload, &body) != nil || body.DurationNS == nil || *body.DurationNS < 0 {
		return "Callback returned: duration unavailable"
	}
	duration := time.Duration(*body.DurationNS).Round(time.Millisecond)
	return fmt.Sprintf("Callback returned: %d ms", duration.Milliseconds())
}

func describeSkip(payload json.RawMessage) string {
	var body struct {
		Reason    string `json:"reason"`
		Condition *struct {
			EntityID      string `json:"entity_id"`
			ExpectedState string `json:"expected_state"`
		} `json:"condition"`
	}
	if json.Unmarshal(payload, &body) != nil {
		return "Skipped: details unavailable"
	}
	text := "Skipped: " + body.Reason
	if body.Condition != nil {
		text += " (" + body.Condition.EntityID + " expected " + body.Condition.ExpectedState + ")"
	}
	return text
}

func describeWait(payload json.RawMessage) string {
	var body struct {
		TargetState string `json:"target_state"`
		DurationNS  int64  `json:"duration_ns"`
	}
	if json.Unmarshal(payload, &body) != nil {
		return "Waiting: details unavailable"
	}
	return "Waiting " + formatDuration(time.Duration(body.DurationNS)) + " for target state " + body.TargetState
}

func describeStateRead(payload json.RawMessage, phase string) string {
	var body struct {
		Method  string `json:"method"`
		Request struct {
			EntityID  string   `json:"entity_id"`
			EntityIDs []string `json:"entity_ids"`
		} `json:"request"`
	}
	if json.Unmarshal(payload, &body) != nil {
		return "State read " + phase + ": details unavailable"
	}
	detail := body.Method
	if body.Request.EntityID != "" {
		detail += " " + body.Request.EntityID
	} else if len(body.Request.EntityIDs) > 0 {
		detail += " " + strings.Join(body.Request.EntityIDs, ", ")
	}
	return "State read " + phase + ": " + detail
}

func describeServiceOperation(payload json.RawMessage, phase string) string {
	var body struct {
		Kind    string `json:"kind"`
		Request struct {
			Domain    string `json:"domain"`
			Service   string `json:"service"`
			EventType string `json:"event_type"`
		} `json:"request"`
	}
	if json.Unmarshal(payload, &body) != nil {
		return "Service operation " + phase + ": details unavailable"
	}
	detail := body.Kind
	if body.Kind == "call" && (body.Request.Domain != "" || body.Request.Service != "") {
		detail = strings.Trim(body.Request.Domain+"."+body.Request.Service, ".")
	} else if body.Kind == "fire_event" && body.Request.EventType != "" {
		detail = body.Request.EventType
	}
	return "Service operation " + phase + ": " + detail
}

func describeSystemEvent(event metrics.StoredEvent) string {
	switch event.Type {
	case metrics.EventSessionStarting:
		return "Session starting"
	case metrics.EventSessionReady:
		return "Session ready"
	case metrics.EventSessionFinished:
		return "Session finished: " + payloadString(event.PayloadJSON, "reason")
	case metrics.EventMetricsDegraded:
		return "Metrics storage degraded"
	case metrics.EventMetricsEventsDropped:
		return "Metrics events dropped"
	case metrics.EventMetricsRecovered:
		return "Metrics storage recovered"
	default:
		return strings.ReplaceAll(string(event.Type), "_", " ")
	}
}

func payloadString(payload json.RawMessage, key string) string {
	var values map[string]any
	if json.Unmarshal(payload, &values) != nil {
		return "details unavailable"
	}
	if value, ok := values[key].(string); ok && value != "" {
		return value
	}
	return "details unavailable"
}

func payloadError(payload json.RawMessage) string {
	var values struct {
		Error *struct {
			GoType  string `json:"go_type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(payload, &values) != nil || values.Error == nil {
		return ""
	}
	if values.Error.GoType == "" {
		return values.Error.Message
	}
	return values.Error.GoType + ": " + values.Error.Message
}

func isOperationFinished(eventType metrics.EventType) bool {
	return eventType == metrics.EventStateReadFinished || eventType == metrics.EventServiceCallFinished
}
