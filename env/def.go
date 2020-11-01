package env

import (
	"errors"
	"go/build"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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

type Paths struct {
	value []string
}

func (p *Paths) Value() interface{} {
	paths := make([]string, 0, len(p.value))
	for _, p := range p.value {
		paths = append(paths, expandPath(p))
	}
	return funk.UniqString(paths)
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

// MarshalYAML implements the interface `yaml.Marshaler`
func (p *Paths) MarshalYAML() (interface{}, error) {
	return p.value, nil
}

// UnmarshalYAML implements the interface `yaml.Unmarshaler`
func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var parsed []string
	if err := unmarshal(&parsed); err != nil {
		return err
	}
	p.value = parsed
	return nil
}

func (p *Paths) MarshalText() (text []byte, err error) {
	return []byte(strings.Join(p.value, string(filepath.ListSeparator))), nil
}

func (p *Paths) UnmarshalText(text []byte) error {
	list := filepath.SplitList(string(text))

	if err := validatePaths(list); err != nil {
		return err
	}
	p.value = funk.UniqString(list)
	return nil
}

func validatePath(path string) (string, error) {
	path = filepath.Clean(path)
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return filepath.EvalSymlinks(path)
	case os.IsNotExist(err):
		return path, nil
	default:
		return "", err
	}
}

func validatePaths(paths []string) error {
	for i, v := range paths {
		r, err := validatePath(v)
		if err != nil {
			return err
		}
		paths[i] = r
	}
	if len(paths) == 0 {
		return errors.New("no path")
	}

	return nil
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
