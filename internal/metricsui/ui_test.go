package metricsui

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"saml.dev/gome-assistant/internal/metrics"
)

func TestActivityHandlerFragmentsExpansionAndRawJSON(t *testing.T) {
	snapshot := Snapshot{SessionID: "session-current", SessionStarted: time.Now().Add(-time.Minute), SessionActive: true, SessionReady: true, Connection: "connected", Running: 1}
	recorder, handler := newTestUI(t, snapshot, []Automation{{
		Key: "kitchen.motion", Type: "entity", Source: "automations.go:42", Callback: "main.motion",
		TriggerSummary: "binary_sensor.kitchen_motion", Configuration: map[string]any{"entity_ids": []string{"binary_sensor.kitchen_motion"}},
	}})

	record(t, recorder, metrics.Event{Type: metrics.EventAutomationTriggered, SessionID: "session-current", AutomationKey: "kitchen.motion", EvaluationID: "eval-1", Payload: metrics.AutomationTriggeredPayload{Trigger: metrics.Trigger{Kind: "entity", Entity: &metrics.EntityData{TriggerEntityID: "binary_sensor.kitchen_motion", FromState: "off", ToState: "on"}}}})
	record(t, recorder, metrics.Event{Type: metrics.EventAutomationCallbackStarted, SessionID: "session-current", AutomationKey: "kitchen.motion", EvaluationID: "eval-1", Payload: metrics.EmptyPayload{}})
	record(t, recorder, metrics.Event{Type: metrics.EventStateReadFinished, SessionID: "session-current", AutomationKey: "kitchen.motion", EvaluationID: "eval-1", OperationID: "operation-1", Payload: metrics.StateReadFinishedPayload{Method: "get", Result: json.RawMessage(`{"trace_only":"large-result-marker"}`), Error: &metrics.Error{GoType: "network.Error", Message: "unavailable"}}})
	record(t, recorder, metrics.Event{Type: metrics.EventAutomationCallbackFinished, SessionID: "session-current", AutomationKey: "kitchen.motion", EvaluationID: "eval-1", Payload: metrics.AutomationCallbackFinishedPayload{DurationNS: int64(time.Second)}})
	record(t, recorder, metrics.Event{Type: metrics.EventSessionReady, SessionID: "session-current", Payload: metrics.EmptyPayload{}})

	response := serve(handler, http.MethodGet, "/?scope=all")
	if response.Code != http.StatusOK {
		t.Fatalf("Activity status = %d; body = %s", response.Code, response.Body.String())
	}
	for _, want := range []string{"Activity", "Automations", "kitchen.motion", "operation warnings", "Session ready", "Session state", "active · ready", `data-uptime="`} {
		if !strings.Contains(response.Body.String(), want) {
			t.Fatalf("Activity body does not contain %q:\n%s", want, response.Body.String())
		}
	}
	if strings.Contains(response.Body.String(), "large-result-marker") {
		t.Fatalf("default Activity response eagerly contains trace payload: %s", response.Body.String())
	}

	automationsOnly := serve(handler, http.MethodGet, "/?scope=automations")
	for _, unwanted := range []string{"Callback started", "State read finished: get", "large-result-marker", "network.Error: unavailable", "Raw JSON"} {
		if strings.Contains(automationsOnly.Body.String(), unwanted) {
			t.Fatalf("initial Activity response unexpectedly contains lazy trace data %q:\n%s", unwanted, automationsOnly.Body.String())
		}
	}

	expanded := serve(handler, http.MethodGet, "/?scope=automations&expand=eval-1")
	for _, want := range []string{`data-evaluation="eval-1"`, `hx-get="/activity/trace?evaluation=eval-1"`, `hx-trigger="load once"`, ` open`, "Loading trace…"} {
		if !strings.Contains(expanded.Body.String(), want) {
			t.Fatalf("expanded Activity shell does not contain %q:\n%s", want, expanded.Body.String())
		}
	}
	if strings.Contains(expanded.Body.String(), "large-result-marker") {
		t.Fatalf("expanded Activity shell eagerly contains trace payload: %s", expanded.Body.String())
	}
	record(t, recorder, metrics.Event{Type: metrics.EventServiceCallFinished, SessionID: "session-current", AutomationKey: "kitchen.motion", EvaluationID: "eval-1", OperationID: "operation-late", Payload: metrics.ServiceCallFinishedPayload{Kind: "call", Result: json.RawMessage(`{"late":"committed-after-list"}`)}})

	trace := serve(handler, http.MethodGet, "/activity/trace?evaluation=eval-1")
	if trace.Code != http.StatusOK {
		t.Fatalf("trace status = %d; body = %s", trace.Code, trace.Body.String())
	}
	for _, want := range []string{"Callback started", "Callback returned: 1000 ms", "State read finished: get", "Raw JSON", "network.Error: unavailable", "large-result-marker", "committed-after-list"} {
		if !strings.Contains(trace.Body.String(), want) {
			t.Fatalf("lazy trace does not contain %q:\n%s", want, trace.Body.String())
		}
	}

	fragment := serve(handler, http.MethodGet, "/activity/fragment?scope=automations&expand=eval-1")
	if got, want := fragment.Header().Get("HX-Push-Url"), "/?expand=eval-1&scope=automations"; got != want {
		t.Fatalf("fragment push URL = %q, want %q", got, want)
	}
	if !strings.Contains(fragment.Body.String(), `id="activity-content"`) || strings.Contains(fragment.Body.String(), "<!doctype") {
		t.Fatalf("Activity fragment has wrong shape: %s", fragment.Body.String())
	}
	if strings.Contains(fragment.Body.String(), "large-result-marker") || strings.Contains(fragment.Body.String(), "committed-after-list") {
		t.Fatalf("Activity fragment eagerly contains trace payload: %s", fragment.Body.String())
	}
}

func TestCallbackTraceRenderingUsesInlineDurationWithoutRawJSON(t *testing.T) {
	started := metrics.StoredEvent{
		Type:        metrics.EventAutomationCallbackStarted,
		PayloadJSON: json.RawMessage(`{}`),
	}
	finished := metrics.StoredEvent{
		Type:        metrics.EventAutomationCallbackFinished,
		PayloadJSON: json.RawMessage(`{"duration_ns":1234567890}`),
	}
	trace := []traceItem{
		{Event: started, Text: describeTraceEvent(started)},
		{Event: finished, Text: describeTraceEvent(finished)},
	}

	rendered := renderEvaluationTrace(trace)
	for _, want := range []string{"Callback started", "Callback returned: 1235 ms"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("callback trace does not contain %q: %s", want, rendered)
		}
	}
	for _, unwanted := range []string{"Raw JSON", "duration_ns"} {
		if strings.Contains(rendered, unwanted) {
			t.Fatalf("callback trace unexpectedly contains %q: %s", unwanted, rendered)
		}
	}
}

func TestActivityTraceEndpointErrorsAreScopedAndVisible(t *testing.T) {
	_, handler := newTestUI(t, Snapshot{}, nil)
	for _, test := range []struct {
		path   string
		status int
		text   string
	}{
		{path: "/activity/trace", status: http.StatusBadRequest, text: "Evaluation ID is required."},
		{path: "/activity/trace?evaluation=missing", status: http.StatusNotFound, text: "This trace is no longer available."},
	} {
		response := serve(handler, http.MethodGet, test.path)
		if response.Code != test.status || !strings.Contains(response.Body.String(), `class="trace-load-error" role="alert"`) || !strings.Contains(response.Body.String(), test.text) {
			t.Fatalf("GET %s = %d, body %q", test.path, response.Code, response.Body.String())
		}
	}
	for _, want := range []string{"htmx:beforeSwap", `path.indexOf('/activity/trace?') === 0`, "event.detail.shouldSwap = true", "path === '/activity/older'", "event.detail.parameters.delete('expand')"} {
		if !strings.Contains(clientScript, want) {
			t.Fatalf("lazy trace client behavior does not contain %q", want)
		}
	}
	retry := renderEvaluationTraceError("Could not load this trace.", "eval & retry")
	for _, want := range []string{`hx-get="/activity/trace?evaluation=eval+%26+retry"`, `hx-target="closest .trace"`, `>Retry</button>`} {
		if !strings.Contains(retry, want) {
			t.Fatalf("trace retry does not contain %q: %s", want, retry)
		}
	}
}

func TestActivityRowsHaveStableScanColumnsAndGroupedResult(t *testing.T) {
	row := evaluationRow{
		ID:            "eval-1",
		AutomationKey: "kitchen.motion",
		Source:        "automations.go:42",
		Trigger:       "binary_sensor.kitchen_motion off → on",
		OccurredAt:    time.Now(),
		Outcome:       "completed",
		Duration:      42 * time.Millisecond,
		Warnings:      2,
		Incomplete:    true,
	}

	rendered := renderActivityRows(activityModel{
		Rows:  []activityRow{{ID: 1, Evaluation: &row}},
		Query: activityQuery{Expanded: map[string]bool{}},
	})
	for _, want := range []string{
		`class="feed-header"`,
		`<span>Time</span><span>Automation</span><span>Trigger</span><span>Result</span>`,
		`hx-get="/activity/trace?evaluation=eval-1" hx-trigger="toggle once" hx-target="find .trace"`,
		`<li class="trace-loading"><span>Loading trace…</span></li>`,
		`<span class="evaluation-meta"><span class="badge completed">completed</span><span class="duration">42ms</span><span class="badge warning-badge">2 operation warnings</span><span class="tooltip"`,
		`</span></summary>`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("Activity rows do not contain %q: %s", want, rendered)
		}
	}
	if !strings.Contains(styles, `--activity-columns: 1.1rem 11rem minmax(11rem, 1fr) minmax(13rem, 1.3fr) 10rem`) {
		t.Fatal("Activity rows do not share fixed time and result tracks")
	}
	system := renderSystem(systemRow{Event: metrics.StoredEvent{OccurredAt: row.OccurredAt, PayloadJSON: json.RawMessage(`{}`)}})
	if strings.Contains(system, `data-show-seconds`) {
		t.Fatalf("system row timestamp is wider than evaluation timestamps: %s", system)
	}
}

func TestMetricsUIUsesFlatSquaredVisualLanguage(t *testing.T) {
	summary := renderActivitySummary(activityModel{})
	filters := renderActivityFilters(activityQuery{})
	automation := renderAutomation(automationRow{Automation: Automation{Key: "daily.one", Type: "daily"}}, false)

	for _, want := range []string{
		`class="summary-actions"`,
		`data-filter-kind="activity"`,
		`class="search-control"`,
		`class="type-chip" data-type="daily"`,
		`background: Canvas; color: CanvasText;`,
		`--wash: color-mix(in srgb, LinkText 7%, Canvas);`,
		`border-radius: 0`,
		`.type-chip[data-type="daily"]`,
		`.type-chip[data-type="interval"]`,
		`.type-chip[data-type="entity"]`,
		`.type-chip[data-type="event"]`,
		`.system-event .raw[open]`,
		`@media (prefers-reduced-motion: reduce)`,
		`@media (forced-colors: active)`,
		`reducedMotion ? 'auto' : 'smooth'`,
	} {
		if !strings.Contains(summary+filters+automation+styles+clientScript, want) {
			t.Fatalf("flat metrics UI does not contain %q", want)
		}
	}
	for _, unwanted := range []string{"box-shadow:", "linear-gradient(", "radial-gradient("} {
		if strings.Contains(styles, unwanted) {
			t.Fatalf("flat metrics UI unexpectedly contains elevated decoration %q", unwanted)
		}
	}
}

func TestAutomationCatalogScrollRegionIsAccessible(t *testing.T) {
	rendered := renderAutomationsContent(automationsModel{
		Rows: []automationRow{{Automation: Automation{Key: "daily.one", Type: "daily"}}},
	})
	for _, want := range []string{
		`class="catalog-table-wrap" role="region" aria-label="Automation catalog" aria-describedby="catalog-scroll-hint" tabindex="0"`,
		`<p id="catalog-scroll-hint" class="scroll-hint">Scroll horizontally to see every column.</p>`,
		`<caption class="sr-only">Current automation definitions and recent timing</caption>`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("accessible automation catalog does not contain %q: %s", want, rendered)
		}
	}
}

func TestInactiveSessionShowsStateWithoutGrowingUptime(t *testing.T) {
	_, handler := newTestUI(t, Snapshot{
		SessionID: "finished", SessionStarted: time.Now().Add(-time.Hour), SessionActive: false, SessionReady: false,
	}, nil)
	response := serve(handler, http.MethodGet, "/")
	if !strings.Contains(response.Body.String(), "inactive") {
		t.Fatalf("inactive session state is missing: %s", response.Body.String())
	}
	if strings.Contains(response.Body.String(), `data-uptime="`) {
		t.Fatalf("inactive session incorrectly renders a live uptime: %s", response.Body.String())
	}
}

func TestAutomationsHandlerAndActivityLink(t *testing.T) {
	recorder, handler := newTestUI(t, Snapshot{SessionID: "s", SessionActive: true}, []Automation{
		{
			Key: "porch.light", Type: "daily", Source: "main.go:9", Callback: "main.porch", TriggerSummary: "sunset + 10m",
			Configuration: map[string]any{"solar_offset": "10m0s", "enabled_conditions": []any{}},
		},
		{Key: "never.run", Type: "daily", Source: "main.go:12", Callback: "main.never", TriggerSummary: "sunrise"},
	})
	firstTrigger := time.Date(2026, time.September, 7, 10, 0, 0, 0, time.UTC)
	completion := firstTrigger.Add(time.Minute)
	lastTrigger := completion.Add(time.Minute)
	record(t, recorder, metrics.Event{Type: metrics.EventAutomationTriggered, OccurredAt: firstTrigger, SessionID: "s", AutomationKey: "porch.light", EvaluationID: "porch-1", Payload: metrics.AutomationTriggeredPayload{Trigger: metrics.Trigger{Kind: "daily"}}})
	record(t, recorder, metrics.Event{Type: metrics.EventAutomationCallbackFinished, OccurredAt: completion, SessionID: "s", AutomationKey: "porch.light", EvaluationID: "porch-1", Payload: metrics.AutomationCallbackFinishedPayload{}})
	record(t, recorder, metrics.Event{Type: metrics.EventAutomationTriggered, OccurredAt: lastTrigger, SessionID: "s", AutomationKey: "porch.light", EvaluationID: "porch-2", Payload: metrics.AutomationTriggeredPayload{Trigger: metrics.Trigger{Kind: "daily"}}})
	record(t, recorder, metrics.Event{Type: metrics.EventAutomationSkipped, OccurredAt: lastTrigger.Add(time.Second), SessionID: "s", AutomationKey: "porch.light", EvaluationID: "porch-2", Payload: metrics.AutomationSkippedPayload{Reason: "condition"}})
	record(t, recorder, metrics.Event{Type: metrics.EventSessionReady, SessionID: "s", Payload: metrics.EmptyPayload{}})

	response := serve(handler, http.MethodGet, "/automations?type=daily&expand=porch.light")
	for _, want := range []string{
		"Current frozen catalog", "porch.light", "sunset + 10m", "Current configuration", "solar_offset",
		`href="/?automation=porch.light"`, "Last trigger", "Last completion",
		`datetime="2026-09-07T10:02:00Z"`, `datetime="2026-09-07T10:01:00Z"`,
		`aria-label="Automation type filters"`, `All automations <strong>2</strong>`,
		`aria-current="page">Daily <strong>2</strong>`, `Intervals <strong>0</strong>`,
		`<form class="filters" action="/automations" method="get" hx-get="/automations/fragment" hx-target="#automations-content" hx-swap="outerHTML"`,
	} {
		if !strings.Contains(response.Body.String(), want) {
			t.Fatalf("Automations body does not contain %q:\n%s", want, response.Body.String())
		}
	}
	for _, unwanted := range []string{"Latest activity", `<span class="badge completed">`, `<span class="badge skipped">`} {
		if strings.Contains(response.Body.String(), unwanted) {
			t.Fatalf("Automations body unexpectedly contains %q:\n%s", unwanted, response.Body.String())
		}
	}
	fragment := serve(handler, http.MethodGet, "/automations/fragment?type=daily")
	if got, want := fragment.Header().Get("HX-Push-Url"), "/automations?type=daily"; got != want {
		t.Fatalf("automations fragment push URL = %q, want %q", got, want)
	}
	activity := serve(handler, http.MethodGet, "/?automation=porch.light")
	if !strings.Contains(activity.Body.String(), "porch.light") || strings.Contains(activity.Body.String(), "Session ready") {
		t.Fatalf("View activity does not strictly apply the automation filter: %s", activity.Body.String())
	}
}

func TestAutomationCountFiltersAreOneClickAndURLBacked(t *testing.T) {
	model := automationsModel{
		Query: automationQuery{
			Search:   "garage door",
			Type:     "event",
			Expanded: map[string]bool{"event.one": true},
		},
		Catalog: []Automation{
			{Key: "daily.one", Type: "daily"},
			{Key: "daily.two", Type: "daily"},
			{Key: "interval.one", Type: "interval"},
			{Key: "entity.one", Type: "entity"},
			{Key: "event.one", Type: "event"},
		},
	}
	rendered := renderAutomationCountFilters(model)

	for _, want := range []string{
		`href="/automations?expand=event.one&amp;q=garage+door"`,
		`hx-get="/automations/fragment?expand=event.one&amp;q=garage+door"`,
		`hx-target="#automations-content" hx-swap="outerHTML"`,
		`href="/automations?expand=event.one&amp;q=garage+door&amp;type=daily"`,
		`All automations <strong>5</strong>`,
		`Daily <strong>2</strong>`,
		`Intervals <strong>1</strong>`,
		`Entity listeners <strong>1</strong>`,
		`aria-current="page">Event listeners <strong>1</strong>`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("automation count filters do not contain %q: %s", want, rendered)
		}
	}
	if got := strings.Count(rendered, `aria-current="page"`); got != 1 {
		t.Fatalf("automation count filters have %d selected actions, want 1: %s", got, rendered)
	}
}

func TestAutomationCatalogSortIsCaseInsensitive(t *testing.T) {
	ui := &ui{options: Options{Automations: func() []Automation {
		return []Automation{
			{Key: "zulu", Type: "event"},
			{Key: "zz-daily", Type: "daily"},
			{Key: "Bravo", Type: "event"},
			{Key: "apple", Type: "event"},
			{Key: "Apricot", Type: "event"},
			{Key: "alpha", Type: "event"},
			{Key: "ALPHA", Type: "event"},
		}
	}}}

	items := ui.catalog()
	keys := make([]string, len(items))
	for index, item := range items {
		keys[index] = item.Key
	}
	want := []string{"zz-daily", "ALPHA", "alpha", "apple", "Apricot", "Bravo", "zulu"}
	if strings.Join(keys, "\n") != strings.Join(want, "\n") {
		t.Fatalf("automation keys = %q, want %q", keys, want)
	}
}

func TestRenderAutomationUsesEmDashForMissingMilestones(t *testing.T) {
	rendered := renderAutomationsContent(automationsModel{
		Rows: []automationRow{{Automation: Automation{Key: "never.run", Type: "event", TriggerSummary: "startup"}}},
	})
	if !strings.Contains(rendered, `<td>startup</td><td class="automation-time">—</td><td class="automation-time">—</td><td class="automation-time">—</td>`) {
		t.Fatalf("missing automation times do not render as em dashes: %s", rendered)
	}
	if strings.Contains(rendered, "Not yet recorded") {
		t.Fatalf("legacy missing-time copy remains: %s", rendered)
	}
}

func TestTimestampsUseCompactFallbackAndFullLocalizedHoverText(t *testing.T) {
	value := time.Date(2026, time.September, 10, 12, 34, 56, 123456789, time.UTC)
	rendered := timestamp(value)
	for _, want := range []string{
		`datetime="2026-09-10T12:34:56.123456789Z"`,
		`title="2026-09-10T12:34:56.123456789Z"`,
		`aria-label="2026-09-10T12:34:56.123456789Z"`,
		`>Sep 10, 12:34 UTC</time>`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("compact timestamp does not contain %q: %s", want, rendered)
		}
	}
	for _, want := range []string{
		`month: 'short'`,
		`second: '2-digit'`,
		`node.hasAttribute('data-show-seconds')`,
		`timeZoneName: 'long'`,
		`node.title = full + ' · ' + node.getAttribute('datetime')`,
		`node.setAttribute('aria-label', full)`,
	} {
		if !strings.Contains(clientScript, want) {
			t.Fatalf("client timestamp localization does not contain %q", want)
		}
	}

	automations := renderAutomationsContent(automationsModel{
		Rows: []automationRow{{Automation: Automation{Key: "timed", Type: "daily", NextDue: value}, LastTriggered: value, LastCompleted: value}},
	})
	if got := strings.Count(automations, `class="automation-time"`); got != 6 {
		t.Fatalf("automation timestamp class count = %d, want 6: %s", got, automations)
	}
	if !strings.Contains(styles, `.catalog .automation-time { width: 1%; white-space: nowrap;`) {
		t.Fatalf("automation timestamp columns are not compact")
	}

	precise := timestampWithSeconds(value)
	for _, want := range []string{`data-show-seconds`, `>Sep 10, 12:34:56 UTC</time>`} {
		if !strings.Contains(precise, want) {
			t.Fatalf("precise compact timestamp does not contain %q: %s", want, precise)
		}
	}
	if !strings.Contains(localTime(value), `aria-hidden="true"`) {
		t.Fatal("relative timestamp is not hidden from assistive technology after the full time label")
	}
}

func TestAutomationCatalogMatchesActivityRowChrome(t *testing.T) {
	for _, want := range []string{
		`.catalog { width: 100%; min-width: 850px; border-collapse: collapse; font-size: .84rem; }`,
		`.catalog tbody tr { cursor: pointer; }`,
		`event.target.closest('.catalog tbody tr')`,
		`details.open = !details.open`,
	} {
		if !strings.Contains(styles+clientScript, want) {
			t.Fatalf("automation catalog row chrome does not contain %q", want)
		}
	}
}

func TestAutomationKeyColumnDoesNotWrap(t *testing.T) {
	rendered := renderAutomationsContent(automationsModel{
		Rows: []automationRow{{Automation: Automation{Key: "a long automation key", Type: "event"}}},
	})
	for _, want := range []string{
		`<th scope="col" class="automation-key">Key</th>`,
		`<td class="automation-key"><details data-automation="a long automation key">`,
		`.catalog th.automation-key, .catalog td.automation-key > details > summary { white-space: nowrap; }`,
	} {
		if !strings.Contains(rendered+styles, want) {
			t.Fatalf("non-wrapping automation key column does not contain %q", want)
		}
	}
}

func TestQueryFiltersAreURLBacked(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?q=kitchen+motion&scope=automations&outcome=completed&type=entity&automation=one&expand=a&expand=b&before=44", nil)
	query := parseActivityQuery(request)
	if query.Search != "kitchen motion" || query.Scope != "automations" || query.Outcome != "completed" || query.AutomationType != "entity" || query.AutomationKey != "one" || query.BeforeID != 44 || !query.Expanded["a"] || !query.Expanded["b"] {
		t.Fatalf("unexpected parsed Activity query: %#v", query)
	}
	url := query.url("/", true)
	for _, want := range []string{"q=kitchen+motion", "scope=automations", "outcome=completed", "type=entity", "automation=one", "before=44", "expand=a", "expand=b"} {
		if !strings.Contains(url, want) {
			t.Fatalf("URL %q does not preserve %q", url, want)
		}
	}
	if got := activityForAutomation("kitchen motion"); got != "/?automation=kitchen+motion" {
		t.Fatalf("View activity URL = %q", got)
	}
}

func TestUnterminatedEvaluationIsOnlyRunningForCurrentActiveSession(t *testing.T) {
	summary := metrics.EvaluationSummary{EvaluationID: "eval", SessionID: "same", AutomationKey: "key", FirstID: 1, LastID: 1, FirstOccurred: time.Now().Add(-time.Second), LastOccurred: time.Now(), TriggerPayload: json.RawMessage(`{"trigger":{"kind":"daily"}}`)}
	definitions := map[string]Automation{"key": {Key: "key", Type: "daily"}}

	finished := summarizeEvaluation(summary, definitions, Snapshot{SessionID: "same", SessionActive: false}, time.Now())
	if finished.Outcome != "incomplete" || !finished.Incomplete {
		t.Fatalf("unterminated finished-session evaluation = %#v, want incomplete", finished)
	}
	active := summarizeEvaluation(summary, definitions, Snapshot{SessionID: "same", SessionActive: true}, time.Now())
	if active.Outcome != "running" || active.Incomplete {
		t.Fatalf("unterminated active-session evaluation = %#v, want running", active)
	}
	older := summarizeEvaluation(summary, definitions, Snapshot{SessionID: "other", SessionActive: true}, time.Now())
	if older.Outcome != "incomplete" {
		t.Fatalf("unterminated older-session evaluation = %#v, want incomplete", older)
	}
	waitingSummary := summary
	waitingSummary.LatestState = metrics.EventAutomationWaiting
	waiting := summarizeEvaluation(waitingSummary, definitions, Snapshot{SessionID: "same", SessionActive: false}, time.Now())
	if waiting.Outcome != "incomplete" {
		t.Fatalf("unterminated waiting evaluation from a finished session = %#v, want incomplete", waiting)
	}
}

func TestHasOlderRequiresOneHundredAndFirstLogicalRows(t *testing.T) {
	recorder, handler := newTestUI(t, Snapshot{}, nil)
	for index := 0; index < activityPageSize; index++ {
		id := fmt.Sprintf("evaluation-%03d", index)
		if !recorder.Record(metrics.Event{Type: metrics.EventAutomationTriggered, SessionID: "finished", AutomationKey: id, EvaluationID: id, Payload: metrics.AutomationTriggeredPayload{Trigger: metrics.Trigger{Kind: "daily"}}}) {
			t.Fatalf("Record(%s) was not admitted", id)
		}
	}
	eventually(t, func() bool {
		summaries, err := recorder.Reader().Evaluations(context.Background(), metrics.EvaluationFilter{Limit: activityScanLimit})
		return err == nil && len(summaries) == activityPageSize
	})
	model, err := handler.(*ui).activity(context.Background(), activityQuery{Scope: "all", Expanded: map[string]bool{}})
	if err != nil {
		t.Fatal(err)
	}
	if model.HasOlder || len(model.Rows) != activityPageSize {
		t.Fatalf("exactly 100 logical rows: HasOlder=%v, rows=%d", model.HasOlder, len(model.Rows))
	}

	if !recorder.Record(metrics.Event{Type: metrics.EventAutomationTriggered, SessionID: "finished", AutomationKey: "evaluation-101", EvaluationID: "evaluation-101", Payload: metrics.AutomationTriggeredPayload{Trigger: metrics.Trigger{Kind: "daily"}}}) {
		t.Fatal("101st evaluation was not admitted")
	}
	eventually(t, func() bool {
		summaries, err := recorder.Reader().Evaluations(context.Background(), metrics.EvaluationFilter{Limit: activityScanLimit})
		return err == nil && len(summaries) == activityPageSize+1
	})
	model, err = handler.(*ui).activity(context.Background(), activityQuery{Scope: "all", Expanded: map[string]bool{}})
	if err != nil {
		t.Fatal(err)
	}
	if !model.HasOlder || len(model.Rows) != activityPageSize {
		t.Fatalf("101 logical rows: HasOlder=%v, rows=%d", model.HasOlder, len(model.Rows))
	}
}

func TestClientClosesEventSourceOnNavigation(t *testing.T) {
	for _, want := range []string{
		"new EventSource('/events')",
		"window.addEventListener('pagehide', closeSSE)",
		"event.persisted) connectSSE()",
		"sse.close()",
	} {
		if !strings.Contains(clientScript, want) {
			t.Fatalf("SSE lifecycle client behavior does not contain %q", want)
		}
	}
	rendered := page("Activity", "activity", "")
	for _, unwanted := range []string{"htmx-sse", "sse-connect", "metrics-sse"} {
		if strings.Contains(rendered, unwanted) {
			t.Fatalf("page still uses the leaking SSE extension: %q in %s", unwanted, rendered)
		}
	}
}

func TestReadOnlyRoutesAndPinnedAssets(t *testing.T) {
	_, handler := newTestUI(t, Snapshot{}, nil)
	for _, path := range []string{"/", "/activity/fragment", "/activity/older", "/activity/trace", "/automations", "/automations/fragment", "/events", "/assets/htmx-2.0.10.min.js"} {
		response := serve(handler, http.MethodPost, path)
		if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
			t.Fatalf("POST %s = %d, Allow %q", path, response.Code, response.Header().Get("Allow"))
		}
	}
	if response := serve(handler, http.MethodGet, "/not-a-route"); response.Code != http.StatusNotFound {
		t.Fatalf("unknown route = %d, want 404", response.Code)
	}
	for _, asset := range []struct{ name, want string }{
		{"assets/htmx-2.0.10.min.js", "71ea67185bfa8c98c39d31717c6fce5d852370fcdfd129db4543774d3145c0de"},
	} {
		data, err := os.ReadFile(asset.name)
		if err != nil {
			t.Fatal(err)
		}
		actual := sha256.Sum256(data)
		if got := hex.EncodeToString(actual[:]); got != asset.want {
			t.Fatalf("%s SHA256 = %s, want %s", asset.name, got, asset.want)
		}
		response := serve(handler, http.MethodGet, "/"+asset.name)
		if response.Code != http.StatusOK || !strings.Contains(response.Header().Get("Content-Type"), "javascript") {
			t.Fatalf("asset %s response = %d, type %q", asset.name, response.Code, response.Header().Get("Content-Type"))
		}
	}
}

func TestNilSourceReturnsServiceUnavailableForDataRoutes(t *testing.T) {
	handler := New(Options{})
	for _, path := range []string{"/", "/activity", "/activity/fragment", "/activity/older", "/activity/trace", "/automations", "/automations/fragment", "/events"} {
		response := serve(handler, http.MethodGet, path)
		if response.Code != http.StatusServiceUnavailable || !strings.Contains(response.Body.String(), errReaderUnavailable.Error()) {
			t.Fatalf("nil source GET %s = %d, body %q", path, response.Code, response.Body.String())
		}
	}
}

func TestSSEInvalidatesOnlyAfterCommittedRecord(t *testing.T) {
	recorder, handler := newTestUI(t, Snapshot{}, nil)
	server := httptest.NewServer(WithCompression(handler))
	defer server.Close()

	requestContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, server.URL+"/events", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Accept-Encoding", "gzip")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.Header.Get("Content-Encoding") != "" {
		t.Fatalf("SSE Content-Encoding = %q, want identity", response.Header.Get("Content-Encoding"))
	}
	reader := bufio.NewReader(response.Body)
	if initial := readSSE(t, reader); initial != "event: change\ndata:\n" {
		t.Fatalf("initial SSE = %q, want payload-free invalidation", initial)
	}
	if !recorder.Record(metrics.Event{Type: metrics.EventSessionReady, SessionID: "s", Payload: metrics.EmptyPayload{}}) {
		t.Fatal("record was not admitted")
	}
	if update := readSSE(t, reader); update != "event: change\ndata:\n" {
		t.Fatalf("SSE update = %q, want payload-free invalidation", update)
	}
	events, err := recorder.Reader().Events(context.Background(), metrics.EventFilter{SessionID: "s"})
	if err != nil || len(events) != 1 {
		t.Fatalf("SSE arrived before committed read was available: events=%#v err=%v", events, err)
	}
}

func newTestUI(t *testing.T, snapshot Snapshot, catalog []Automation) (*metrics.Recorder, http.Handler) {
	t.Helper()
	recorder, err := metrics.OpenRecorder(t.TempDir() + "/metrics.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = recorder.Close(nil) })
	return recorder, New(Options{
		Source:      recorder,
		Snapshot:    func() Snapshot { return snapshot },
		Automations: func() []Automation { return catalog },
		Now:         time.Now,
	})
}

func record(t *testing.T, recorder *metrics.Recorder, event metrics.Event) {
	t.Helper()
	changes, cancel := recorder.SubscribeChanges()
	defer cancel()
	if !recorder.Record(event) {
		t.Fatalf("Record(%s) was not admitted", event.Type)
	}
	select {
	case <-changes:
	case <-time.After(3 * time.Second):
		t.Fatalf("Record(%s) did not commit", event.Type)
	}
}

func eventually(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition did not become true")
}

func serve(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func readSSE(t *testing.T, reader *bufio.Reader) string {
	t.Helper()
	result := make(chan struct {
		value string
		err   error
	}, 1)
	go func() {
		var lines strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				result <- struct {
					value string
					err   error
				}{err: err}
				return
			}
			if line == "\n" {
				result <- struct {
					value string
					err   error
				}{value: lines.String()}
				return
			}
			lines.WriteString(line)
		}
	}()
	select {
	case read := <-result:
		if read.err != nil && read.err != io.EOF {
			t.Fatal(read.err)
		}
		return read.value
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for SSE")
		return ""
	}
}
