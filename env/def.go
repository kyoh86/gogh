package env

import (
	"go/build"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kyoh86/appenv/types"
	"github.com/kyoh86/gogh/gogh"
	"github.com/thoas/go-funk"
)

const (
	KeyringService = "gogh.kyoh86.dev"
	EnvarPrefix    = "GOGH_"
)

type GithubHost struct {
	types.StringPropertyBase
}

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost = "github.com"
)

func (*GithubHost) Default() interface{} {
	return DefaultHost
}

type GithubUser struct {
	types.StringPropertyBase
}

type Roots struct {
	value []string
}

func (p *Roots) Value() interface{} {
	var list []string
	uniq := map[string]struct{}{}
	for _, p := range p.value {
		exp := expandPath(p)
		if _, ok := uniq[exp]; ok {
			continue
		}
		uniq[exp] = struct{}{}
		list = append(list, exp)
	}
	return list
}

func expandPath(path string) string {
	if len(path) == 0 {
		return path
	}

	path = os.ExpandEnv(path)
	if path[0] != '~' || (len(path) > 1 && path[1] != filepath.Separator) {
		return path
	}

	user, err := user.Current()
	if err != nil {
		return path
	}

	return filepath.Join(user.HomeDir, path[1:])
}
func (*Roots) Default() interface{} {
	gopaths := filepath.SplitList(build.Default.GOPATH)
	root := make([]string, 0, len(gopaths))
	for _, gopath := range gopaths {
		root = append(root, filepath.Join(gopath, "src"))
	}
	return funk.UniqString(root)
}

// MarshalYAML implements the interface `yaml.Marshaler`
func (p *Roots) MarshalYAML() (interface{}, error) {
	return p.value, nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (p *Roots) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed []string
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	p.value = parsed
	return nil
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
