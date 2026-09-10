package metricsui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func renderAutomationsContent(model automationsModel) string {
	var out strings.Builder
	out.WriteString(`<div id="automations-content"><section class="catalog-heading"><h1>Automations</h1><p class="muted">Current frozen catalog. Definitions shown here are not historical snapshots.</p></section>`)
	out.WriteString(renderAutomationCountFilters(model))
	out.WriteString(`<form class="filters" action="/automations" method="get" hx-get="/automations/fragment" hx-target="#automations-content" hx-swap="outerHTML" hx-push-url="true" data-filter-kind="automations"><label class="search-control">Search <input type="search" name="q" value="` + escape(model.Query.Search) + `" placeholder="key, source, trigger"></label>`)
	out.WriteString(selectControl("type", "Automation type", model.Query.Type, []option{{"", "All types"}, {"daily", "Daily"}, {"interval", "Interval"}, {"entity", "Entity"}, {"event", "Event"}}))
	out.WriteString(`<button type="submit">Filter</button></form>`)
	if len(model.Rows) == 0 {
		out.WriteString(`<p class="empty">No current automations match these filters.</p></div>`)
		return out.String()
	}
	out.WriteString(`<div class="catalog-table-wrap" role="region" aria-label="Automation catalog" aria-describedby="catalog-scroll-hint" tabindex="0"><p id="catalog-scroll-hint" class="scroll-hint">Scroll horizontally to see every column.</p><table class="catalog"><caption class="sr-only">Current automation definitions and recent timing</caption><thead><tr><th scope="col" class="automation-key">Key</th><th scope="col">Type</th><th scope="col">Source</th><th scope="col">Trigger</th><th scope="col" class="automation-time">Last trigger</th><th scope="col" class="automation-time">Last completion</th><th scope="col" class="automation-time">Next due</th></tr></thead><tbody>`)
	for _, row := range model.Rows {
		out.WriteString(renderAutomation(row, model.Query.Expanded[row.Key]))
	}
	out.WriteString(`</tbody></table></div></div>`)
	return out.String()
}

func renderAutomationCountFilters(model automationsModel) string {
	var out strings.Builder
	out.WriteString(`<nav class="summary-controls automation-counts" aria-label="Automation type filters">`)
	types := append([]struct{ typeName, label string }{{"", "All automations"}}, automationKinds...)
	for _, kind := range types {
		query := model.Query
		query.Type = kind.typeName
		count := len(model.Catalog)
		if kind.typeName != "" {
			count = countAutomations(model.Catalog, kind.typeName)
		}
		current := ""
		if model.Query.Type == kind.typeName {
			current = ` aria-current="page"`
		}
		out.WriteString(`<a href="` + escape(query.url("/automations")) + `" hx-get="` + escape(query.url("/automations/fragment")) + `" hx-target="#automations-content" hx-swap="outerHTML" hx-push-url="true"` + current + `>` + escape(kind.label) + ` <strong>` + fmt.Sprintf("%d", count) + `</strong></a>`)
	}
	out.WriteString(`</nav>`)
	return out.String()
}

func renderAutomation(row automationRow, expanded bool) string {
	open := ""
	if expanded {
		open = " open"
	}
	lastTriggered := automationTime(row.LastTriggered)
	lastCompleted := automationTime(row.LastCompleted)
	nextDue := "—"
	if !row.NextDue.IsZero() {
		nextDue = localTime(row.NextDue)
	}
	return `<tr><td class="automation-key"><details data-automation="` + escape(row.Key) + `"` + open + `><summary><strong>` + escape(row.Key) + `</strong></summary><div class="automation-details"><p><strong>Callback:</strong> ` + escape(row.Callback) + `</p><p><strong>Current configuration</strong></p><pre>` + escape(jsonConfigurationOrMessage(row.Configuration)) + `</pre><p><a href="` + escape(activityForAutomation(row.Key)) + `">View activity</a></p></div></details></td><td><span class="type-chip" data-type="` + escape(row.Type) + `">` + escape(row.Type) + `</span></td><td>` + escape(row.Source) + `</td><td>` + escape(row.TriggerSummary) + `</td><td class="automation-time">` + lastTriggered + `</td><td class="automation-time">` + lastCompleted + `</td><td class="automation-time">` + nextDue + `</td></tr>`
}

func automationTime(value time.Time) string {
	if value.IsZero() {
		return "—"
	}
	return localTime(value)
}

func jsonConfiguration(values map[string]any) (string, error) {
	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func jsonConfigurationOrMessage(values map[string]any) string {
	value, err := jsonConfiguration(values)
	if err != nil {
		return fmt.Sprintf("Configuration is unavailable: %v", err)
	}
	return value
}

func activityForAutomation(key string) string {
	return activityQuery{AutomationKey: key}.url("/", false)
}
