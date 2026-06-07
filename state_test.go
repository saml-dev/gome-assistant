package gomeassistant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatesEmpty(t *testing.T) {
	state := StateImpl{}

	states, err := state.GetStates(nil)

	assert.NoError(t, err)
	assert.Equal(t, []string{}, states)
}
