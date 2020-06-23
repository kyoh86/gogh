package command_test

import (
	"testing"

	"github.com/kyoh86/gogh/command"
	"github.com/stretchr/testify/assert"
)

func ExampleInitialize_zsh() {
	if err := command.Initialize(nil, "gogh-cd", "zsh"); err != nil {
		panic(err)
	}
	// Output:
	// function gogh() {
	//   exec 5>&1
	//   case $1 in
	//   "cd" )
	//     shift
	//     cd "$(command gogh find "$@" | tee /dev/tty | tail -n1)" || return
	//     ;;
	//
	//   "get" )
	//     local CD=0
	//     for arg in "$@"; do
	//       if [ "${arg}" = '--cd' ]; then
	//         CD=1
	//       fi
	//     done
	//
	//     if [ $CD -eq 1 ]; then
	//       loc="$(command gogh "$@" | tee /dev/tty | tail -n1)"
	//       cd "$loc" || return
	//     else
	//       command gogh "$@"
	//     fi
	//     ;;
	//
	//   * )
	//     command gogh "$@"
	//     ;;
	//   esac
	// }
	// eval "$(command gogh --completion-script-zsh)"
}
func ExampleInitialize_bash() {
	if err := command.Initialize(nil, "gogh-cd", "bash"); err != nil {
		panic(err)
	}
	// Output:
	// #!/bin/bash
	//
	// function gogh() {
	//   case $1 in
	//   "cd" )
	//     shift
	//     cd "$(command gogh find "$@" | tee /dev/tty | tail -n1)" || return
	//     ;;
	//
	//   "get" )
	//     local CD=0
	//     for arg in "$@"; do
	//       if [ "${arg}" = '--cd' ]; then
	//         CD=1
	//       fi
	//     done
	//
	//     if [ $CD -eq 1 ]; then
	//       loc="$(command gogh "$@" | tee /dev/tty | tail -n1)"
	//       cd "$loc" || return
	//     else
	//       command gogh "$@"
	//     fi
	//     ;;
	//
	//   * )
	//     command gogh "$@"
	//     ;;
	//   esac
	// }
	// eval "$(command gogh --completion-script-bash)"
}

func TestInitialize(t *testing.T) {
	assert.EqualError(t, command.Initialize(nil, "gogh-cd", "invalid"), "unsupported shell \"invalid\"")
	assert.NoError(t, command.Initialize(nil, "gogh-cd", "zsh"))
	assert.NoError(t, command.Initialize(nil, "gogh-cd", "bash"))
}
