package config

import (
	"errors"
	"strings"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/util"
)

type ConfigAccessor map[string]OptionAccessor

var (
	DefaultConfigAccessor = NewConfigAccessor()
)

func (m ConfigAccessor) Accessor(optionName string) (*OptionAccessor, error) {
	a, ok := m[optionName]
	if !ok {
		return nil, InvalidOptionName
	}
	return &a, nil
}

func (m ConfigAccessor) OptionNames() []string {
	arr := make([]string, 0, len(m))
	for o := range m {
		arr = append(arr, o)
	}
	return arr
}

type OptionAccessor struct {
	optionName string
	getter     func(ctx gogh.Context) string
	setter     func(config *Config, value string) error
	unsetter   func(config *Config) error
}

func (a OptionAccessor) Get(ctx gogh.Context) string            { return a.getter(ctx) }
func (a OptionAccessor) Set(config *Config, value string) error { return a.setter(config, value) }
func (a OptionAccessor) Unset(config *Config) error             { return a.unsetter(config) }

var (
	EmptyValue           = errors.New("empty value")
	RemoveFromMonoOption = errors.New("removing from mono option")
	InvalidOptionName    = errors.New("invalid option name")
)

// TODO: generate

func NewConfigAccessor() ConfigAccessor {
	m := ConfigAccessor{}
	for _, a := range []OptionAccessor{
		GitHubUserOptionAccessor,
		GitHubTokenOptionAccessor,
		GitHubHostOptionAccessor,
		LogLevelOptionAccessor,
		RootOptionAccessor,
	} {
		m[a.optionName] = a
	}
	return m
}

var (
	GitHubUserOptionAccessor = OptionAccessor{
		optionName: "github.user",
		getter: func(ctx gogh.Context) string {
			return ctx.GitHubUser()
		},
		setter: func(config *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			if err := gogh.ValidateOwner(value); err != nil {
				return err
			}
			config.GitHub.User = value
			return nil
		},
		unsetter: func(config *Config) error {
			config.GitHub.User = ""
			return nil
		},
	}

	GitHubTokenOptionAccessor = OptionAccessor{
		optionName: "github.token",
		getter: func(ctx gogh.Context) string {
			return ctx.GitHubToken()
		},
		setter: func(config *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			config.GitHub.Token = value
			return nil
		},
		unsetter: func(config *Config) error {
			config.GitHub.Token = ""
			return nil
		},
	}

	GitHubHostOptionAccessor = OptionAccessor{
		optionName: "github.host",
		getter: func(ctx gogh.Context) string {
			return ctx.GitHubHost()
		},
		setter: func(config *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			config.GitHub.Host = value
			return nil
		},
		unsetter: func(config *Config) error {
			config.GitHub.Host = ""
			return nil
		},
	}

	LogLevelOptionAccessor = OptionAccessor{
		optionName: "loglevel",
		getter: func(ctx gogh.Context) string {
			return ctx.LogLevel()
		},
		setter: func(config *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			if err := gogh.ValidateLogLevel(value); err != nil {
				return err
			}
			config.VLogLevel = value
			return nil
		},
		unsetter: func(config *Config) error {
			config.VLogLevel = ""
			return nil
		},
	}

	RootOptionAccessor = OptionAccessor{
		optionName: "root",
		getter: func(ctx gogh.Context) string {
			return strings.Join(ctx.Root(), "\n")
		},
		setter: func(config *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			path, err := gogh.ValidateRoot(value)
			if err != nil {
				return err
			}
			config.VRoot = util.UniqueStringArray(append(config.VRoot, path))
			return nil
		},
		unsetter: func(config *Config) error {
			config.VRoot = nil
			return nil
		},
	}
)
