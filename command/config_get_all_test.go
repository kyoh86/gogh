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
  - /root1
  - /root2
hooks:
  - /hook1
  - /hook2
githubHost: hostx1
githubUser: userx1`)
	config, err := env.GetConfig(yml)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigGetAll(&config); err != nil {
		log.Fatalln(err)
	}

	// Unordered output:
	// roots: /root1:/root2
	// hooks: /hook1:/hook2
	// github.host: hostx1
	// github.user: userx1
	// github.token: *****
}

func TestConfigGetAll(t *testing.T) {
	yml := strings.NewReader(`
roots:
  - /root1
  - /root2
hooks:
  - /hook1
  - /hook2
githubHost: hostx1`)
	config, err := env.GetConfig(yml)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigGetAll(&config))
}
