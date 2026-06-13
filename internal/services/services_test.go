package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTargetSerializesMultipleEntityIDs(t *testing.T) {
	req := BaseServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_off",
		Target:  Entities([]string{"light.kitchen", "fan.office"}),
	}

	body, err := json.Marshal(req)

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"domain": "homeassistant",
		"service": "turn_off",
		"target": {
			"entity_id": ["light.kitchen", "fan.office"]
		}
	}`, string(body))
}

func TestTargetSerializesSingleEntityIDList(t *testing.T) {
	target := Entities([]string{"light.kitchen"})

	body, err := json.Marshal(target)

	require.NoError(t, err)
	assert.JSONEq(t, `{
		"entity_id": ["light.kitchen"]
	}`, string(body))
}
