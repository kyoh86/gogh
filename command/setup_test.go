package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
	"github.com/stretchr/testify/assert"
)

func ExampleSetupForZsh() {
	if err := command.Setup(&config.Config{}, "gogh-cd", "zsh"); err != nil {
		panic(err)
	}
	// Output:
	// function gogh-cd { cd $(gogh find $@) }
	// eval "$(gogh --completion-script-zsh)"
}

func ExampleSetupForBash() {
	if err := command.Setup(&config.Config{}, "gogh-cd", "bash"); err != nil {
		panic(err)
	}
	// Output:
	// function gogh-cd { cd $(gogh find $@) }
	// eval "$(gogh --completion-script-bash)"
}

func TestSetup(t *testing.T) {
	assert.EqualError(t, command.Setup(&config.Config{}, "gogh-cd", "invalid"), "unsupported shell \"invalid\"")
}
