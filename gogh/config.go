package gogh

import (
	"context"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joeshaw/envdecode"
	"github.com/pelletier/go-toml"
)

// Config holds configuration file values.
type Config struct {
	context.Context
	VLogLevel string       `toml:"loglevel,omitempty" env:"GOGH_LOG_LEVEL"`
	VRoot     RootConfig   `toml:"root,omitempty" env:"GOGH_ROOT"`
	GitHub    GitHubConfig `toml:"github,omitempty"`
}

func (c *Config) Stdout() io.Writer {
	return os.Stdout
}

func (c *Config) Stderr() io.Writer {
	return os.Stderr
}

func (c *Config) GitHubUser() string {
	return c.GitHub.User
}

func (c *Config) GitHubToken() string {
	return c.GitHub.Token
}

func (c *Config) GitHubHost() string {
	return c.GitHub.Host
}

func (c *Config) LogLevel() string {
	return c.VLogLevel
}

func (c *Config) Root() []string {
	return c.VRoot
}

func (c *Config) PrimaryRoot() string {
	return c.VRoot[0]
}

type RootConfig []string

// Decode implements the interface `envdecode.Decoder`
func (r *RootConfig) Decode(repl string) error {
	*r = strings.Split(repl, ":")
	return nil
}

type GitHubConfig struct {
	Token string `toml:"token,omitempty" env:"GOGH_GITHUB_TOKEN"`
	User  string `toml:"user,omitempty" env:"GOGH_GITHUB_USER"`
	Host  string `toml:"host,omitempty" env:"GOGH_GITHUB_HOST"`
}

var (
	envGoghLogLevel    = "GOGH_LOG_LEVEL"
	envGoghGitHubUser  = "GOGH_GITHUB_USER"
	envGoghGitHubToken = "GOGH_GITHUB_TOKEN"
	envGoghGitHubHost  = "GOGH_GITHUB_HOST"
	envGoghRoot        = "GOGH_ROOT"
	envNames           = []string{
		envGoghLogLevel,
		envGoghGitHubUser,
		envGoghGitHubToken,
		envGoghGitHubHost,
		envGoghRoot,
	}
)

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost     = "github.com"
	DefaultLogLevel = "warn"
)

var defaultConfig = Config{
	VLogLevel: DefaultLogLevel,
	GitHub: GitHubConfig{
		Host: DefaultHost,
	},
}

var initDefaultConfig sync.Once

func DefaultConfig() *Config {
	initDefaultConfig.Do(func() {
		gopaths := filepath.SplitList(build.Default.GOPATH)
		root := make([]string, 0, len(gopaths))
		for _, gopath := range gopaths {
			root = append(root, filepath.Join(gopath, "src"))
		}
		defaultConfig.VRoot = unique(root)
	})
	return &defaultConfig
}

func LoadFileConfig(filename string) (config *Config, err error) {
	file, err := os.Open(filename)
	switch {
	case err == nil:
		defer file.Close()
		config = &Config{}
		err = toml.NewDecoder(file).Decode(config)
	case os.IsNotExist(err):
		err = nil
	}
	config.VRoot = unique(config.VRoot)
	return
}

func GetEnvarConfig() (config *Config, err error) {
	config = &Config{}
	err = envdecode.Decode(config)
	if err == envdecode.ErrNoTargetFieldsAreSet {
		err = nil
	}
	config.VRoot = unique(config.VRoot)
	return
}

func MergeConfig(base *Config, override ...*Config) *Config {
	c := base
	for _, o := range override {
		c.VLogLevel = mergeStringOption(c.VLogLevel, o.VLogLevel)
		c.VRoot = mergeStringArrayOption(c.VRoot, o.VRoot)
		c.GitHub.Token = mergeStringOption(c.GitHub.Token, o.GitHub.Token)
		c.GitHub.User = mergeStringOption(c.GitHub.User, o.GitHub.User)
		c.GitHub.Host = mergeStringOption(c.GitHub.Host, o.GitHub.Host)
	}
	return c
}

func mergeStringOption(base, override string) string {
	if override != "" {
		return override
	}
	return base
}

func mergeStringArrayOption(base, override []string) []string {
	if len(override) > 0 {
		return override
	}
	return base
}
