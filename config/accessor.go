package config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/gogh"
	"github.com/thoas/go-funk"
	"github.com/zalando/go-keyring"
)

var (
	ErrEmptyValue        = errors.New("empty value")
	ErrInvalidOptionName = errors.New("invalid option name")
)

type OptionAccessor struct {
	optionName string
	getter     func(cfg *Config) string
	putter     func(cfg *Config, value string) error
	unsetter   func(cfg *Config) error
}

func (a OptionAccessor) Get(cfg *Config) string              { return a.getter(cfg) }
func (a OptionAccessor) Put(cfg *Config, value string) error { return a.putter(cfg, value) }
func (a OptionAccessor) Unset(cfg *Config) error             { return a.unsetter(cfg) }

var (
	configAccessor  map[string]OptionAccessor
	optionNames     []string
	optionAccessors = []OptionAccessor{
		rootOptionAccessor,
		gitHubHostOptionAccessor,
		gitHubUserOptionAccessor,
		gitHubTokenOptionAccessor,
	}
)

func init() {
	m := map[string]OptionAccessor{}
	n := make([]string, 0, len(optionAccessors))
	for _, a := range optionAccessors {
		n = append(n, a.optionName)
		m[a.optionName] = a
	}
	configAccessor = m
	optionNames = n
}

func Option(optionName string) (*OptionAccessor, error) {
	a, ok := configAccessor[optionName]
	if !ok {
		return nil, ErrInvalidOptionName
	}
	return &a, nil
}

func OptionNames() []string {
	return optionNames
}

var (
	gitHubUserOptionAccessor = OptionAccessor{
		optionName: "github.user",
		getter: func(cfg *Config) string {
			return cfg.GithubUser()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := gogh.ValidateOwner(value); err != nil {
				return err
			}
			cfg.GitHub.User = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.User = ""
			return nil
		},
	}

	gitHubTokenOptionAccessor = OptionAccessor{
		optionName: "github.token",
		getter: func(cfg *Config) string {
			if cfg.GithubToken() == "" {
				return ""
			}
			return "*****"
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return keyring.Set(keyGoghServiceName, keyGoghGithubToken, value)
		},
		unsetter: func(cfg *Config) error {
			return keyring.Delete(keyGoghServiceName, keyGoghGithubToken)
		},
	}

	gitHubHostOptionAccessor = OptionAccessor{
		optionName: "github.host",
		getter: func(cfg *Config) string {
			return cfg.GithubHost()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			cfg.GitHub.Host = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.Host = ""
			return nil
		},
	}

	rootOptionAccessor = OptionAccessor{
		optionName: "root",
		getter: func(cfg *Config) string {
			return strings.Join(cfg.Root(), string(filepath.ListSeparator))
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}

			list := filepath.SplitList(value)

			if err := gogh.ValidateRoots(list); err != nil {
				return err
			}
			cfg.VRoot = funk.UniqString(append(cfg.VRoot, list...))
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VRoot = nil
			return nil
		},
	}
)
