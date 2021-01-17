package config

import (
	"go/build"
	"path/filepath"

	"github.com/kyoh86/appenv/types"
	"github.com/kyoh86/xdg"
	"github.com/thoas/go-funk"
)

const (
	KeyringService = "gogh.kyoh86.dev"
	EnvarPrefix    = "GOGH_"
)

type GithubHost struct {
	types.StringValue
}

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost = "github.com"
)

func (*GithubHost) Default() interface{} {
	return DefaultHost
}

type GithubUser struct {
	types.StringValue
}

type Roots struct {
	Paths
}

func (*Roots) Default() interface{} {
	gopaths := filepath.SplitList(build.Default.GOPATH)
	paths := make([]string, 0, len(gopaths))
	for _, gopath := range gopaths {
		paths = append(paths, filepath.Join(gopath, "src"))
	}
	return funk.UniqString(paths)
}

type Hooks struct {
	Paths
}

func (*Hooks) Default() interface{} {
	return []string{filepath.Join(xdg.ConfigHome(), "gogh", "hooks")}
}
