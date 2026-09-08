package metricsui

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	activityPageSize  = 100
	activityScanLimit = 1000
)

type activityQuery struct {
	Search         string
	Scope          string
	Outcome        string
	AutomationType string
	AutomationKey  string
	BeforeID       int64
	Expanded       map[string]bool
}

func parseActivityQuery(request *http.Request) activityQuery {
	values := request.URL.Query()
	query := activityQuery{
		Search:         trimTo(values.Get("q"), 256),
		Scope:          allowed(values.Get("scope"), "all", "automations", "system"),
		Outcome:        allowed(values.Get("outcome"), "", "waiting", "running", "completed", "skipped", "canceled", "incomplete"),
		AutomationType: allowed(values.Get("type"), "", "daily", "interval", "entity", "event"),
		AutomationKey:  trimTo(values.Get("automation"), 256),
		Expanded:       make(map[string]bool),
	}
	if query.Scope == "" {
		query.Scope = "all"
	}
	if before, err := strconv.ParseInt(values.Get("before"), 10, 64); err == nil && before > 0 {
		query.BeforeID = before
	}
	for _, id := range values["expand"] {
		if id = trimTo(id, 128); id != "" && len(query.Expanded) < 50 {
			query.Expanded[id] = true
		}
	}
	return query
}

func (query activityQuery) values(includeBefore bool) url.Values {
	values := make(url.Values)
	if query.Search != "" {
		values.Set("q", query.Search)
	}
	if query.Scope != "" && query.Scope != "all" {
		values.Set("scope", query.Scope)
	}
	if query.Outcome != "" {
		values.Set("outcome", query.Outcome)
	}
	if query.AutomationType != "" {
		values.Set("type", query.AutomationType)
	}
	if query.AutomationKey != "" {
		values.Set("automation", query.AutomationKey)
	}
	if includeBefore && query.BeforeID > 0 {
		values.Set("before", strconv.FormatInt(query.BeforeID, 10))
	}
	for id := range query.Expanded {
		values.Add("expand", id)
	}
	return values
}

func (query activityQuery) url(path string, includeBefore bool) string {
	values := query.values(includeBefore)
	if encoded := values.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

type automationQuery struct {
	Search   string
	Type     string
	Expanded map[string]bool
}

func parseAutomationQuery(request *http.Request) automationQuery {
	values := request.URL.Query()
	query := automationQuery{
		Search:   trimTo(values.Get("q"), 256),
		Type:     allowed(values.Get("type"), "", "daily", "interval", "entity", "event"),
		Expanded: make(map[string]bool),
	}
	for _, key := range values["expand"] {
		if key = trimTo(key, 128); key != "" && len(query.Expanded) < 50 {
			query.Expanded[key] = true
		}
	}
	return query
}

func (query automationQuery) url(path string) string {
	values := make(url.Values)
	if query.Search != "" {
		values.Set("q", query.Search)
	}
	if query.Type != "" {
		values.Set("type", query.Type)
	}
	for key := range query.Expanded {
		values.Add("expand", key)
	}
	if encoded := values.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

func allowed(value string, values ...string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	for _, candidate := range values {
		if value == candidate {
			return value
		}
	}
	return ""
}

func trimTo(value string, maximum int) string {
	value = strings.TrimSpace(value)
	if len(value) > maximum {
		return value[:maximum]
	}
	return value
}
