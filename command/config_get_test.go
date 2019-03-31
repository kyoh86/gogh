package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigGet() {
	if err := command.ConfigGet(&config.Config{
		VRoot: []string{"/foo", "/bar"},
	}, "root"); err != nil {
		panic(err)
	}
	// Output:
	// /foo:/bar
}

func TestConfigGet(t *testing.T) {
	assert.EqualError(t, command.ConfigGet(&config.Config{}, "invalid.name"), "invalid option name")
}
