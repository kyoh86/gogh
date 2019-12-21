package sh

import (
	"time"

	"github.com/deadcheat/goblet"
)

//go:generate goblet -g --ignore-dotfiles -o assets.go -p sh src

// Assets a generated file system
var Assets = goblet.NewFS(
	map[string][]string{
		"/src": {
			"init.bash", "init.zsh",
		},
	},
	map[string]*goblet.File{
		"/src/init.bash": goblet.NewFile("/src/init.bash", _Assetsa96e64fd7a6ce6e1579ce8d99622cfdcf8c12b96, 0x1a4, time.Unix(1575764391, 1575764391416360881)),
		"/src/init.zsh":  goblet.NewFile("/src/init.zsh", _Assetsd275f00e2d9fab456599903fda40166d6da9aa46, 0x1a4, time.Unix(1575764381, 1575764381889652313)),
		"/src":           goblet.NewFile("/src", nil, 0x800001ed, time.Unix(1575764391, 1575764391416360881)),
	},
)

// binary data
var (
	_Assetsa96e64fd7a6ce6e1579ce8d99622cfdcf8c12b96 = []byte("#!/bin/bash\n\nfunction gogh() {\n  case $1 in\n  \"cd\" )\n    shift\n    cd \"$(command gogh find \"$@\" | tee /dev/tty | tail -n1)\" || return\n    ;;\n\n  \"get\" )\n    local CD=0\n    for arg in \"$@\"; do\n      if [ \"${arg}\" = '--cd' ]; then\n        CD=1\n      fi\n    done\n\n    if [ $CD -eq 1 ]; then\n      loc=\"$(command gogh \"$@\" | tee /dev/tty | tail -n1)\"\n      cd \"$loc\" || return\n    else\n      command gogh \"$@\"\n    fi\n    ;;\n\n  * )\n    command gogh \"$@\"\n    ;;\n  esac\n}\neval \"$(command gogh --completion-script-bash)\"\n")
	_Assetsd275f00e2d9fab456599903fda40166d6da9aa46 = []byte("function gogh() {\n  exec 5>&1\n  case $1 in\n  \"cd\" )\n    shift\n    cd \"$(command gogh find \"$@\" | tee /dev/tty | tail -n1)\" || return\n    ;;\n\n  \"get\" )\n    local CD=0\n    for arg in \"$@\"; do\n      if [ \"${arg}\" = '--cd' ]; then\n        CD=1\n      fi\n    done\n\n    if [ $CD -eq 1 ]; then\n      loc=\"$(command gogh \"$@\" | tee /dev/tty | tail -n1)\"\n      cd \"$loc\" || return\n    else\n      command gogh \"$@\"\n    fi\n    ;;\n\n  * )\n    command gogh \"$@\"\n    ;;\n  esac\n}\neval \"$(command gogh --completion-script-zsh)\"\n")
)
