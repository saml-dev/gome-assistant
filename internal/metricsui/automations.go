package metricsui

import (
	"net/http"
	"strings"
	"time"

	"saml.dev/gome-assistant/internal/metrics"
)

type automationsModel struct {
	Query   automationQuery
	Catalog []Automation
	Rows    []automationRow
}

type automationRow struct {
	Automation
	LastTriggered time.Time
	LastCompleted time.Time
}

func (ui *ui) automations(request *http.Request) (automationsModel, error) {
	query := parseAutomationQuery(request)
	catalog := ui.catalog()
	model := automationsModel{Query: query, Catalog: catalog}
	for _, item := range catalog {
		if query.Type != "" && item.Type != query.Type {
			continue
		}
		if query.Search != "" && !matchesAutomation(item, query.Search) {
			continue
		}
		model.Rows = append(model.Rows, automationRow{Automation: item})
	}
	if ui.options.Source == nil || ui.options.Source.Reader() == nil {
		return model, errReaderUnavailable
	}
	keys := make([]string, 0, len(model.Rows))
	for _, row := range model.Rows {
		keys = append(keys, row.Key)
	}
	milestones, err := ui.options.Source.Reader().LatestAutomationMilestones(request.Context(), keys)
	if err != nil {
		return model, err
	}
	latest := make(map[string]metrics.AutomationMilestones, len(milestones))
	for _, milestone := range milestones {
		latest[milestone.AutomationKey] = milestone
	}
	for index := range model.Rows {
		milestone := latest[model.Rows[index].Key]
		model.Rows[index].LastTriggered = milestone.LastTriggered
		model.Rows[index].LastCompleted = milestone.LastCompleted
	}
	return model, nil
}

func matchesAutomation(item Automation, search string) bool {
	configuration, _ := jsonConfiguration(item.Configuration)
	return containsFold(strings.Join([]string{item.Key, item.Type, item.Source, item.Callback, item.TriggerSummary, configuration}, " "), search)
}
