package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigPut(t *testing.T) {
	var cfg config.Config
	assert.NoError(t, command.ConfigPut(&cfg, "github.host", "hostx1"))
	assert.EqualError(t, command.ConfigPut(&cfg, "invalid.name", "hostx2"), "invalid option name")
}
