package main

import (
	"testing"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
}

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	i, e := checkArgs(event)
	assert.Equal(sensu.CheckStateWarning, i)
	assert.Error(e)
	plugin.CriticalMultiplier = float64(2)
	i, e = checkArgs(event)
	assert.Equal(sensu.CheckStateWarning, i)
	assert.Error(e)
	plugin.WarningMultiplier = float64(1.5)
	i, e = checkArgs(event)
	assert.Equal(sensu.CheckStateOK, i)
	assert.NoError(e)
	plugin.WarningMultiplier = float64(3)
	i, e = checkArgs(event)
	assert.Equal(sensu.CheckStateWarning, i)
	assert.Error(e)
}
