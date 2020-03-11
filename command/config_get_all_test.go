package command_test

import (
	"log"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/env"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigGetAll() {
	yml := strings.NewReader(`
roots:
  - /foo
  - /bar
githubHost: hostx1`)
	config, err := env.GetConfig(yml, "")
	if err != nil {
		log.Fatal(err)
	}
	if err := command.ConfigGetAll(&config); err != nil {
		log.Fatal(err)
	}

	// Unordered output:
	// roots: /foo:/bar
	// github.host: hostx1
	// github.token:
}

func TestConfigGetAll(t *testing.T) {
	yml := strings.NewReader(`
roots:
  - /foo
  - /bar
githubHost: hostx1`)
	config, err := env.GetConfig(yml, "")
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigGetAll(&config))
}
