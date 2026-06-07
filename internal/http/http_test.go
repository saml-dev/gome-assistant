package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatesTemplate(t *testing.T) {
	template, err := statesTemplate([]string{"light.kitchen", "sensor.temperature"})

	assert.NoError(t, err)
	assert.Equal(t, `{{ [states("light.kitchen"), states("sensor.temperature")] | to_json }}`, template)
}

func TestStatesTemplateQuotesEntityIDs(t *testing.T) {
	template, err := statesTemplate([]string{`sensor."quoted"`})

	assert.NoError(t, err)
	assert.Equal(t, `{{ [states("sensor.\"quoted\"")] | to_json }}`, template)
}
