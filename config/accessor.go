package config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/util"
)

type Accessor map[string]OptionAccessor

var (
	DefaultAccessor = NewAccessor()
)

func (m Accessor) Option(optionName string) (*OptionAccessor, error) {
	a, ok := m[optionName]
	if !ok {
		return nil, InvalidOptionName
	}
	return &a, nil
}

func (m Accessor) OptionNames() []string {
	arr := make([]string, 0, len(m))
	for o := range m {
		arr = append(arr, o)
	}
	return arr
}

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
	EmptyValue           = errors.New("empty value")
	RemoveFromMonoOption = errors.New("removing from mono option")
	InvalidOptionName    = errors.New("invalid option name")
	TokenMustNotSave     = errors.New("token must not save")
)

func NewAccessor() Accessor {
	m := Accessor{}
	for _, a := range []OptionAccessor{
		GitHubUserOptionAccessor,
		GitHubTokenOptionAccessor,
		GitHubHostOptionAccessor,
		LogLevelOptionAccessor,
		LogDateOptionAccessor,
		LogTimeOptionAccessor,
		LogLongFileOptionAccessor,
		LogShortFileOptionAccessor,
		LogUTCOptionAccessor,
		RootOptionAccessor,
	} {
		m[a.optionName] = a
	}
	return m
}

var (
	GitHubUserOptionAccessor = OptionAccessor{
		optionName: "github.user",
		getter: func(cfg *Config) string {
			return cfg.GitHubUser()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
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

	GitHubTokenOptionAccessor = OptionAccessor{
		optionName: "github.token",
		getter: func(cfg *Config) string {
			return cfg.GitHubToken()
		},
		putter: func(cfg *Config, value string) error {
			return TokenMustNotSave
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.Token = ""
			return nil
		},
	}

	GitHubHostOptionAccessor = OptionAccessor{
		optionName: "github.host",
		getter: func(cfg *Config) string {
			return cfg.GitHubHost()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			cfg.GitHub.Host = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.Host = ""
			return nil
		},
	}

	LogLevelOptionAccessor = OptionAccessor{
		optionName: "log.level",
		getter: func(cfg *Config) string {
			return cfg.LogLevel()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			if err := gogh.ValidateLogLevel(value); err != nil {
				return err
			}
			cfg.Log.Level = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Level = ""
			return nil
		},
	}

	LogDateOptionAccessor = OptionAccessor{
		optionName: "log.date",
		getter: func(cfg *Config) string {
			return cfg.Log.Date.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			return cfg.Log.Date.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Date = EmptyBoolOption
			return nil
		},
	}

	LogTimeOptionAccessor = OptionAccessor{
		optionName: "log.time",
		getter: func(cfg *Config) string {
			return cfg.Log.Time.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			return cfg.Log.Time.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Time = EmptyBoolOption
			return nil
		},
	}

	LogLongFileOptionAccessor = OptionAccessor{
		optionName: "log.longfile",
		getter: func(cfg *Config) string {
			return cfg.Log.LongFile.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			return cfg.Log.LongFile.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.LongFile = EmptyBoolOption
			return nil
		},
	}

	LogShortFileOptionAccessor = OptionAccessor{
		optionName: "log.shortfile",
		getter: func(cfg *Config) string {
			return cfg.Log.ShortFile.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			return cfg.Log.ShortFile.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.ShortFile = EmptyBoolOption
			return nil
		},
	}

	LogUTCOptionAccessor = OptionAccessor{
		optionName: "log.utc",
		getter: func(cfg *Config) string {
			return cfg.Log.UTC.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}
			return cfg.Log.UTC.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.UTC = EmptyBoolOption
			return nil
		},
	}

	RootOptionAccessor = OptionAccessor{
		optionName: "root",
		getter: func(cfg *Config) string {
			return strings.Join(cfg.Root(), string(filepath.ListSeparator))
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return EmptyValue
			}

			list := filepath.SplitList(value)

			if err := gogh.ValidateRoots(list); err != nil {
				return err
			}
			cfg.VRoot = util.UniqueStringArray(append(cfg.VRoot, list...))
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VRoot = nil
			return nil
		},
	}
)
