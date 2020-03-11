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
githubHost: hostx1`)
	config, err := env.GetConfig(source, "")
	if err != nil {
		log.Fatal(err)
	}
	if err := command.ConfigSet(&config, "github.host", "hostx2"); err != nil {
		log.Fatal(err)
	}
	if err := config.Save(os.Stdout, env.DiscardKeyringService); err != nil {
		log.Fatal(err)
	}
	if err := command.ConfigGetAll(&config); err != nil {
		log.Fatal(err)
	}

	// Unordered output:
	// roots:
	//   - /foo
	// githubHost: hostx2
	// roots: /foo
	// github.host: hostx2
	// github.token:
}

func TestConfigSet(t *testing.T) {
	source := strings.NewReader(`
roots:
  - /foo
githubHost: hostx1`)
	config, err := env.GetConfig(source, "")
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigSet(&config, "github.host", "hostx2"))
	assert.NoError(t, config.Save(os.Stdout, env.DiscardKeyringService))
	assert.NoError(t, command.ConfigGetAll(&config))
}
