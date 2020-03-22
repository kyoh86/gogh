package command_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/env"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigSet() {
	source := strings.NewReader(`
roots:
  - /root1
hooks:
  - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigSet(&access, &config, "github.host", "hostx2"); err != nil {
		log.Fatalln(err)
	}
	if err := config.Save(os.Stdout); err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigGetAll(&config); err != nil {
		log.Fatalln(err)
	}

	// Unordered output:
	// roots:
	//   - /root1
	// hooks:
	//   - /hook1
	// githubHost: hostx2
	// githubUser: userx1
	// roots: /root1
	// hooks: /hook1
	// github.host: hostx2
	// github.user: userx1
	// github.token: *****
}

func TestConfigSet(t *testing.T) {
	// NOTE: never use real host name. github.token breaks keyring store
	source := strings.NewReader(`
roots:
  - /root1
hooks:
  - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigSet(&access, &config, "github.host", "hostx2"))
	assert.NoError(t, config.Save(os.Stdout))
	assert.NoError(t, command.ConfigGetAll(&config))

	assert.Error(t, command.ConfigSet(&access, &config, "invalid.config", "invalid"))
	assert.Error(t, command.ConfigUnset(&access, &config, "github.token"))
	assert.NoError(t, command.ConfigSet(&access, &config, "github.token", "invalid"))
	assert.NoError(t, command.ConfigUnset(&access, &config, "github.token"))
}
