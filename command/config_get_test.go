package command_test

import (
	"log"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/env"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigGet() {
	yml := strings.NewReader(`{"roots": ["/foo", "/bar"]}`)
	config, err := env.GetConfig(yml)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigGet(&config, "roots"); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// /foo:/bar
}

func TestConfigGet(t *testing.T) {
	config, err := env.GetConfig(env.EmptyYAMLReader)
	if err != nil {
		log.Fatalln(err)
	}
	assert.EqualError(t, command.ConfigGet(&config, "invalid.name"), `invalid property name "invalid.name"`)
}
