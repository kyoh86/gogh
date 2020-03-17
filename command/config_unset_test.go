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
  - /foo
githubUser: userx1
githubHost: hostx1`)
	config, err := env.GetConfig(source)
	if err != nil {
		log.Fatal(err)
	}
	if err := command.ConfigUnset(&config, "github.host"); err != nil {
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
	// githubUser: userx1
	// roots: /foo
	// github.host:
	// github.user: userx1
	// github.token: *****
}

func TestConfigUnset(t *testing.T) {
	source := strings.NewReader(`
roots:
  - /foo
githubUser: userx1
githubHost: hostx1`)
	config, err := env.GetConfig(source)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigUnset(&config, "github.host"))
	assert.NoError(t, config.Save(os.Stdout))
	assert.NoError(t, command.ConfigGetAll(&config))
	assert.Error(t, command.ConfigUnset(&config, "invalid.config"))
}
