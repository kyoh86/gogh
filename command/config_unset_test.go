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

func ExampleConfigUnset() {
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
	if err := command.ConfigUnset(&access, &config, "github.host"); err != nil {
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
	// githubUser: userx1
	// roots: /root1
	// hooks: /hook1
	// github.host:
	// github.user: userx1
	// github.token: *****
}

func TestConfigUnset(t *testing.T) {
	source := strings.NewReader(`
roots:
  - /root1
hooks:
  - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigUnset(&access, &config, "github.host"))
	assert.NoError(t, config.Save(os.Stdout))
	assert.NoError(t, command.ConfigGetAll(&config))
	assert.Error(t, command.ConfigUnset(&access, &config, "invalid.config"))
}
