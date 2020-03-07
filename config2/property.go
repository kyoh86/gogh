package config

import (
	"go/build"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/config2/types"
	"github.com/kyoh86/gogh/gogh"
	"github.com/thoas/go-funk"
)

type GithubToken struct {
	types.StringPropertyBase
}

func (p *GithubToken) StoreKeyring() {}
func (p *GithubToken) StoreEnvar()   {}

func (p *GithubToken) Mask(value string) string {
	if value == "" {
		return ""
	}
	return "*****"
}

type GithubHost struct {
	types.StringPropertyBase
}

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost = "github.com"
)

func (*GithubHost) StoreConfigFile() {}
func (*GithubHost) StoreEnvar()      {}

func (*GithubHost) Default() interface{} {
	return DefaultHost
}

type GithubUser struct {
	types.StringPropertyBase
}

func (*GithubUser) StoreCacheFile() {}

type Roots struct {
	value []string
}

func (*Roots) StoreConfigFile() {}
func (*Roots) StoreEnvar()      {}

func (p *Roots) Value() interface{} {
	return p.value
}
func (*Roots) Default() interface{} {
	gopaths := filepath.SplitList(build.Default.GOPATH)
	root := make([]string, 0, len(gopaths))
	for _, gopath := range gopaths {
		root = append(root, filepath.Join(gopath, "src"))
	}
	return funk.UniqString(root)
}

func (p *Roots) MarshalText() (text []byte, err error) {
	return []byte(strings.Join(p.value, string(filepath.ListSeparator))), nil
}

func (p *Roots) UnmarshalText(text []byte) error {
	list := filepath.SplitList(string(text))

	if err := gogh.ValidateRoots(list); err != nil {
		return err
	}
	p.value = funk.UniqString(list)
	return nil
}
