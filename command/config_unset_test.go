package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigUnset(t *testing.T) {
	cfg := config.Config{
		GitHub: config.GitHubConfig{
			Host: "hostx1",
		},
	}
	assert.NoError(t, command.ConfigUnset(&cfg, "github.host"))
	assert.Empty(t, cfg.GitHub.Host)
	assert.EqualError(t, command.ConfigUnset(&cfg, "invalid.name"), "invalid option name")
}
