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
  - /foo
githubUser: userx1
githubHost: hostx1`)
	config, err := env.GetConfig(source)
	if err != nil {
		log.Fatal(err)
	}
	if err := command.ConfigSet(&config, "github.host", "hostx2"); err != nil {
		log.Fatal(err)
	}
	if err := config.Save(os.Stdout); err != nil {
		log.Fatal(err)
	}
	if err := command.ConfigGetAll(&config); err != nil {
		log.Fatal(err)
	}

	// Unordered output:
	// roots:
	//   - /foo
	// githubHost: hostx2
	// githubUser: userx1
	// roots: /foo
	// github.host: hostx2
	// github.user: userx1
	// github.token: *****
}

func TestConfigSet(t *testing.T) {
	// NOTE: never use real host name. github.token breaks keyring store
	source := strings.NewReader(`
roots:
  - /foo
githubUser: userx1
githubHost: hostx1`)
	config, err := env.GetConfig(source)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigSet(&config, "github.host", "hostx2"))
	assert.NoError(t, config.Save(os.Stdout))
	assert.NoError(t, command.ConfigGetAll(&config))

	assert.Error(t, command.ConfigSet(&config, "invalid.config", "invalid"))
	assert.Error(t, command.ConfigUnset(&config, "github.token"))
	assert.NoError(t, command.ConfigSet(&config, "github.token", "invalid"))
	assert.NoError(t, command.ConfigUnset(&config, "github.token"))
}
