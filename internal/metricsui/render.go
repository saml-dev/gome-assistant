package metricsui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"

	"saml.dev/gome-assistant/internal/metrics"
)

func (ui *ui) handleActivity(writer http.ResponseWriter, request *http.Request) {
	model, err := ui.activity(request.Context(), parseActivityQuery(request))
	if err != nil {
		ui.writePageError(writer, "Activity", err)
		return
	}
	writeHTML(writer, http.StatusOK, page("Activity", "activity", renderActivityContent(model)))
}

func (ui *ui) handleActivityFragment(writer http.ResponseWriter, request *http.Request) {
	model, err := ui.activity(request.Context(), parseActivityQuery(request))
	if err != nil {
		writeHTML(writer, http.StatusInternalServerError, fragmentError("Could not refresh activity", err))
		return
	}
	writer.Header().Set("HX-Push-Url", model.Query.url("/", false))
	writeHTML(writer, http.StatusOK, renderActivityContent(model))
}

func (ui *ui) handleOlder(writer http.ResponseWriter, request *http.Request) {
	model, err := ui.activity(request.Context(), parseActivityQuery(request))
	if err != nil {
		writeHTML(writer, http.StatusInternalServerError, `<p id="load-older" class="error" role="alert">Could not load older activity.</p>`)
		return
	}
	writeHTML(writer, http.StatusOK, renderActivityRows(model))
}

func (ui *ui) handleActivityTrace(writer http.ResponseWriter, request *http.Request) {
	evaluationID := trimTo(request.URL.Query().Get("evaluation"), 128)
	if evaluationID == "" {
		writeHTML(writer, http.StatusBadRequest, renderEvaluationTraceError("Evaluation ID is required.", ""))
		return
	}
	trace, err := ui.evaluationTrace(request.Context(), evaluationID)
	if err != nil {
		writeHTML(writer, http.StatusInternalServerError, renderEvaluationTraceError("Could not load this trace.", evaluationID))
		return
	}
	if len(trace) == 0 {
		writeHTML(writer, http.StatusNotFound, renderEvaluationTraceError("This trace is no longer available.", ""))
		return
	}
	writeHTML(writer, http.StatusOK, renderEvaluationTrace(trace))
}

func (ui *ui) handleAutomations(writer http.ResponseWriter, request *http.Request) {
	model, err := ui.automations(request)
	if err != nil {
		ui.writePageError(writer, "Automations", err)
		return
	}
	writeHTML(writer, http.StatusOK, page("Automations", "automations", renderAutomationsContent(model)))
}

func (ui *ui) handleAutomationsFragment(writer http.ResponseWriter, request *http.Request) {
	model, err := ui.automations(request)
	if err != nil {
		writeHTML(writer, http.StatusInternalServerError, fragmentError("Could not refresh automations", err))
		return
	}
	writer.Header().Set("HX-Push-Url", model.Query.url("/automations"))
	writeHTML(writer, http.StatusOK, renderAutomationsContent(model))
}

func (ui *ui) writePageError(writer http.ResponseWriter, title string, err error) {
	writeHTML(writer, http.StatusInternalServerError, page(title, "", fragmentError("Metrics UI is unavailable", err)))
}

func writeHTML(writer http.ResponseWriter, status int, body string) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(body))
}

func page(title, active, content string) string {
	activityClass, automationClass := "", ""
	if active == "activity" {
		activityClass = ` class="active" aria-current="page"`
	}
	if active == "automations" {
		automationClass = ` class="active" aria-current="page"`
	}
	return `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="color-scheme" content="light dark"><title>` + escape(title) + ` · gomeassistant metrics</title>
<style>` + styles + `</style>
<script defer src="/assets/htmx-2.0.10.min.js"></script>
<script defer src="/assets/htmx-sse-2.2.4.js"></script></head>
<body><header class="topbar"><nav aria-label="Primary"><a href="/"` + activityClass + `>Activity</a><a href="/automations"` + automationClass + `>Automations</a></nav><span class="product">gomeassistant metrics</span></header>
<div id="metrics-sse" hx-ext="sse" sse-connect="/events" sse-swap="change" aria-hidden="true"></div>
<main>` + content + `</main><script>` + clientScript + `</script></body></html>`
}

func renderActivityContent(model activityModel) string {
	var out strings.Builder
	out.WriteString(`<div id="activity-content">`)
	out.WriteString(renderActivitySummary(model))
	out.WriteString(renderActivityFilters(model.Query))
	out.WriteString(`<p class="new-updates" data-updates-banner hidden aria-live="polite"><button type="button" data-show-new>Show <span data-updates-count>0</span> new updates</button></p>`)
	out.WriteString(renderActivityRows(model))
	out.WriteString(`</div>`)
	return out.String()
}

func renderActivitySummary(model activityModel) string {
	var out strings.Builder
	snapshot := model.Snapshot
	if snapshot.MetricsHealth.Degraded || lossCount(snapshot.MetricsHealth.PendingLosses) > 0 {
		out.WriteString(`<aside class="warning" role="status"><strong>Telemetry may be incomplete.</strong> SQLite is degraded or events are awaiting a loss summary; new activity can be missing.</aside>`)
	}
	out.WriteString(`<section class="summary" aria-label="Current metrics summary"><div><h1>Activity</h1><p class="muted">Committed automation evaluations and system events.</p></div><dl>`)
	out.WriteString(summaryTerm("Metrics", healthText(snapshot.MetricsHealth)))
	connection := snapshot.Connection
	if connection == "" {
		connection = "disconnected"
	}
	out.WriteString(summaryTerm("Home Assistant", connection))
	if !snapshot.SessionStarted.IsZero() {
		out.WriteString(`<div><dt>Session started</dt><dd>` + timestamp(snapshot.SessionStarted))
		if snapshot.SessionActive {
			out.WriteString(`<span class="relative" data-uptime="` + fmt.Sprintf("%d", snapshot.SessionStarted.UnixMilli()) + `"></span>`)
		}
		out.WriteString(`</dd></div>`)
	} else {
		out.WriteString(summaryTerm("Session started", "No session has started"))
	}
	out.WriteString(summaryTerm("Session state", sessionState(snapshot)))
	if snapshot.LatestError != nil {
		out.WriteString(summaryTerm("Latest error", snapshot.LatestError.Message))
	}
	out.WriteString(`</dl><div class="summary-actions"><div class="summary-controls" aria-label="Automation counts">`)
	for _, kind := range automationKinds {
		out.WriteString(`<a href="/automations?type=` + kind.typeName + `">` + escape(kind.label) + ` <strong>` + fmt.Sprintf("%d", countAutomations(model.Catalog, kind.typeName)) + `</strong></a>`)
	}
	out.WriteString(`</div><div class="summary-controls" aria-label="Active evaluation counts">`)
	running := model.Query
	running.Outcome = "running"
	running.Scope = "automations"
	running.BeforeID = 0
	waiting := model.Query
	waiting.Outcome = "waiting"
	waiting.Scope = "automations"
	waiting.BeforeID = 0
	out.WriteString(`<a href="` + escape(running.url("/", false)) + `">Running <strong>` + fmt.Sprintf("%d", snapshot.Running) + `</strong></a>`)
	out.WriteString(`<a href="` + escape(waiting.url("/", false)) + `">Waiting <strong>` + fmt.Sprintf("%d", snapshot.Waiting) + `</strong></a></div></div></section>`)
	return out.String()
}

func renderActivityFilters(query activityQuery) string {
	var out strings.Builder
	out.WriteString(`<form class="filters" action="/" method="get" hx-get="/activity/fragment" hx-target="#activity-content" hx-push-url="true" data-filter-kind="activity"><label class="search-control">Search <input type="search" name="q" value="` + escape(query.Search) + `" placeholder="key, source, entity, event"></label>`)
	if query.AutomationKey != "" {
		out.WriteString(`<input type="hidden" name="automation" value="` + escape(query.AutomationKey) + `"><span class="filter-note">Automation: ` + escape(query.AutomationKey) + ` <a href="` + escape(clearAutomationURL(query)) + `">clear</a></span>`)
	}
	out.WriteString(selectControl("scope", "Scope", query.Scope, []option{{"all", "All"}, {"automations", "Automations"}, {"system", "System"}}))
	out.WriteString(selectControl("outcome", "Outcome", query.Outcome, []option{{"", "Any outcome"}, {"waiting", "Waiting"}, {"running", "Running"}, {"completed", "Completed"}, {"skipped", "Skipped"}, {"canceled", "Canceled"}, {"incomplete", "Incomplete"}}))
	out.WriteString(selectControl("type", "Automation type", query.AutomationType, []option{{"", "Any type"}, {"daily", "Daily"}, {"interval", "Interval"}, {"entity", "Entity"}, {"event", "Event"}}))
	out.WriteString(`<button type="submit">Filter</button></form><nav class="scope" aria-label="Activity scope">`)
	for _, item := range []struct{ value, label string }{{"all", "All"}, {"automations", "Automations"}, {"system", "System"}} {
		link := query
		link.Scope = item.value
		link.BeforeID = 0
		current := ""
		if query.Scope == item.value {
			current = ` aria-current="page"`
		}
		out.WriteString(`<a href="` + escape(link.url("/", false)) + `"` + current + `>` + item.label + `</a>`)
	}
	out.WriteString(`</nav>`)
	return out.String()
}

func renderActivityRows(model activityModel) string {
	var out strings.Builder
	if len(model.Rows) == 0 {
		if model.Query.BeforeID > 0 {
			out.WriteString(`<p class="empty" id="load-older">No older matching activity.</p>`)
			return out.String()
		}
		out.WriteString(`<section id="activity-feed" class="feed"><p class="empty">No activity has been recorded yet. The current <a href="/automations">automation catalog</a> is still available.</p></section>`)
		return out.String()
	}
	if model.Query.BeforeID == 0 {
		out.WriteString(`<section id="activity-feed" class="feed" aria-label="Activity feed"><div class="feed-header" aria-hidden="true"><span></span><span>Time</span><span>Automation</span><span>Trigger</span><span>Result</span></div>`)
	}
	oldest := model.Rows[len(model.Rows)-1].ID
	for _, row := range model.Rows {
		if row.Evaluation != nil {
			out.WriteString(renderEvaluation(*row.Evaluation, model.Query.Expanded[row.Evaluation.ID]))
		} else if row.System != nil {
			out.WriteString(renderSystem(*row.System))
		}
	}
	if model.HasOlder {
		next := model.Query
		next.BeforeID = oldest
		out.WriteString(`<a id="load-older" class="load-older" href="` + escape(next.url("/", true)) + `" hx-get="` + escape(next.url("/activity/older", true)) + `" hx-target="#load-older" hx-swap="outerHTML">Load older</a>`)
	} else {
		out.WriteString(`<p id="load-older" class="end">No older matching activity.</p>`)
	}
	if model.Query.BeforeID == 0 {
		out.WriteString(`</section>`)
	}
	return out.String()
}

func renderEvaluation(row evaluationRow, expanded bool) string {
	var out strings.Builder
	open := ""
	trigger := "toggle once"
	if expanded {
		open = " open"
		trigger = "load once"
	}
	out.WriteString(`<details class="evaluation" data-evaluation="` + escape(row.ID) + `" hx-get="` + escape(evaluationTraceURL(row.ID)) + `" hx-trigger="` + trigger + `" hx-target="find .trace" hx-swap="outerHTML"` + open + `><summary><span class="time">` + localTime(row.OccurredAt) + `</span><span class="evaluation-main"><strong>` + escape(row.AutomationKey) + `</strong><small>` + escape(row.Source) + `</small></span><span class="trigger">` + escape(row.Trigger) + `</span><span class="evaluation-meta"><span class="badge ` + escape(row.Outcome) + `">` + escape(row.Outcome) + `</span><span class="duration">` + escape(formatDuration(row.Duration)) + `</span>`)
	if row.Warnings > 0 {
		out.WriteString(`<span class="badge warning-badge">` + fmt.Sprintf("%d", row.Warnings) + ` operation warnings</span>`)
	}
	if row.Incomplete {
		out.WriteString(`<span class="tooltip" tabindex="0" aria-label="No terminal telemetry event was recorded. The UI cannot distinguish a process crash, callback panic, or dropped event." title="No terminal telemetry event was recorded. The UI cannot distinguish a process crash, callback panic, or dropped event.">incomplete?</span>`)
	}
	out.WriteString(`</span></summary><ol class="trace" aria-live="polite"><li class="trace-loading"><span>Loading trace…</span></li></ol></details>`)
	return out.String()
}

func renderEvaluationTrace(trace []traceItem) string {
	var out strings.Builder
	out.WriteString(`<ol class="trace" aria-live="polite">`)
	for _, item := range trace {
		out.WriteString(`<li>` + timestampWithSeconds(item.Event.OccurredAt) + `<span>` + escape(item.Text) + `</span>`)
		if item.Error != "" {
			out.WriteString(`<strong class="trace-error">` + escape(item.Error) + `</strong>`)
		}
		out.WriteString(traceRawDisclosure(item.Event))
		out.WriteString(`</li>`)
	}
	out.WriteString(`</ol>`)
	return out.String()
}

func renderEvaluationTraceError(message, retryEvaluationID string) string {
	var out strings.Builder
	out.WriteString(`<ol class="trace" aria-live="polite"><li class="trace-load-error" role="alert"><strong>` + escape(message) + `</strong>`)
	if retryEvaluationID != "" {
		out.WriteString(`<button type="button" hx-get="` + escape(evaluationTraceURL(retryEvaluationID)) + `" hx-target="closest .trace" hx-swap="outerHTML">Retry</button>`)
	}
	out.WriteString(`</li></ol>`)
	return out.String()
}

func evaluationTraceURL(evaluationID string) string {
	values := url.Values{"evaluation": []string{evaluationID}}
	return "/activity/trace?" + values.Encode()
}

func traceRawDisclosure(event metrics.StoredEvent) string {
	switch event.Type {
	case metrics.EventAutomationCallbackStarted, metrics.EventAutomationCallbackFinished:
		return ""
	default:
		return rawDisclosure(event.PayloadJSON)
	}
}

func renderSystem(row systemRow) string {
	return `<article class="system-event">` + timestamp(row.Event.OccurredAt) + `<strong>System</strong><span>` + escape(row.Text) + `</span>` + rawDisclosure(row.Event.PayloadJSON) + `</article>`
}

func rawDisclosure(payload json.RawMessage) string {
	return `<details class="raw"><summary>Raw JSON</summary><pre>` + escape(prettyJSON(payload)) + `</pre></details>`
}

func prettyJSON(payload json.RawMessage) string {
	var out bytes.Buffer
	if json.Indent(&out, payload, "", "  ") == nil {
		return out.String()
	}
	return string(payload)
}

func localTime(value time.Time) string {
	return timestamp(value) + `<span class="relative" data-relative="` + fmt.Sprintf("%d", value.UnixMilli()) + `" aria-hidden="true"></span>`
}

func timestamp(value time.Time) string {
	return renderTimestamp(value, false)
}

func timestampWithSeconds(value time.Time) string {
	return renderTimestamp(value, true)
}

func renderTimestamp(value time.Time, withSeconds bool) string {
	utc := value.UTC()
	full := utc.Format(time.RFC3339Nano)
	secondsAttribute := ""
	fallbackFormat := "Jan 2, 15:04 UTC"
	if withSeconds {
		secondsAttribute = ` data-show-seconds`
		fallbackFormat = "Jan 2, 15:04:05 UTC"
	}
	return `<time datetime="` + escape(full) + `" data-local-time` + secondsAttribute + ` title="` + escape(full) + `" aria-label="` + escape(full) + `">` + escape(utc.Format(fallbackFormat)) + `</time>`
}

func summaryTerm(term, value string) string {
	return `<div><dt>` + escape(term) + `</dt><dd>` + escape(value) + `</dd></div>`
}

func healthText(health metrics.Health) string {
	if health.Closed {
		return "closed"
	}
	if health.Degraded {
		return "degraded"
	}
	return "recording"
}

func sessionState(snapshot Snapshot) string {
	if !snapshot.SessionActive {
		return "inactive"
	}
	if snapshot.SessionReady {
		return "active · ready"
	}
	return "active · starting"
}

func lossCount(counts metrics.LossCounts) uint64 {
	return counts.QueueFull + counts.DatabaseError + counts.SerializationError
}

func countAutomations(items []Automation, kind string) int {
	count := 0
	for _, item := range items {
		if item.Type == kind {
			count++
		}
	}
	return count
}

var automationKinds = []struct{ typeName, label string }{
	{"daily", "Daily"},
	{"interval", "Intervals"},
	{"entity", "Entity listeners"},
	{"event", "Event listeners"},
}

type option struct{ value, label string }

func selectControl(name, label, selected string, options []option) string {
	var out strings.Builder
	out.WriteString(`<label>` + escape(label) + ` <select name="` + escape(name) + `">`)
	for _, option := range options {
		isSelected := ""
		if selected == option.value {
			isSelected = " selected"
		}
		out.WriteString(`<option value="` + escape(option.value) + `"` + isSelected + `>` + escape(option.label) + `</option>`)
	}
	out.WriteString(`</select></label>`)
	return out.String()
}

func clearAutomationURL(query activityQuery) string {
	query.AutomationKey = ""
	return query.url("/", false)
}

func escape(value string) string { return html.EscapeString(value) }

func formatDuration(value time.Duration) string {
	if value < 0 {
		value = 0
	}
	if value < time.Second {
		return value.Round(time.Millisecond).String()
	}
	return value.Round(time.Millisecond).String()
}

func fragmentError(message string, err error) string {
	return `<section class="error" role="alert"><strong>` + escape(message) + `.</strong> ` + escape(err.Error()) + `</section>`
}
