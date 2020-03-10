package config

import (
	"context"
	"time"
)

// Config holds configuration file values.
type Config struct {
	VRoot  PathListOption `yaml:"root,omitempty" env:"GOGH_ROOT"`
	GitHub GitHubConfig   `yaml:"github,omitempty"`
}

var _ context.Context = (*Config)(nil)

type GitHubConfig struct {
	Token string `yaml:"-" env:"GOGH_GITHUB_TOKEN"`
	User  string `yaml:"user,omitempty" env:"GOGH_GITHUB_USER"`
	Host  string `yaml:"host,omitempty" env:"GOGH_GITHUB_HOST"`
}

// Deadline : empty context.Context
func (*Config) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done : empty context.Context
func (*Config) Done() <-chan struct{} {
	return nil
}

// Err : empty context.Context
func (*Config) Err() error {
	return nil
}

// Value : empty context.Context
func (*Config) Value(key interface{}) interface{} {
	return nil
}

func (*Config) String() string {
	return "gogh/config.Config"
}

func (c *Config) GithubUser() string {
	return c.GitHub.User
}

func (c *Config) GithubToken() string {
	return c.GitHub.Token
}

func (c *Config) GithubHost() string {
	return c.GitHub.Host
}

func (c *Config) Root() []string {
	ret := make([]string, 0, len(c.VRoot))
	for _, p := range c.VRoot {
		ret = append(ret, expandPath(p))
	}
	return ret
}

func (c *Config) PrimaryRoot() string {
	return expandPath(c.VRoot[0])
}
