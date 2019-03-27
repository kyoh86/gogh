package config

import (
	"go/build"
	"io"
	"path/filepath"
	"sync"

	"github.com/joeshaw/envdecode"
	"github.com/kyoh86/gogh/internal/util"
	yaml "gopkg.in/yaml.v2"
)

var (
	envGoghLogLevel        = "GOGH_LOG_LEVEL"
	envGoghLogDate         = "GOGH_LOG_DATE"
	envGoghLogTime         = "GOGH_LOG_TIME"
	envGoghLogMicroSeconds = "GOGH_LOG_MICROSECONDS"
	envGoghLogLongFile     = "GOGH_LOG_LONGFILE"
	envGoghLogShortFile    = "GOGH_LOG_SHORTFILE"
	envGoghLogUTC          = "GOGH_LOG_UTC"
	envGoghGitHubUser      = "GOGH_GITHUB_USER"
	envGoghGitHubToken     = "GOGH_GITHUB_TOKEN"
	envGoghGitHubHost      = "GOGH_GITHUB_HOST"
	envGoghRoot            = "GOGH_ROOT"
	envNames               = []string{
		envGoghLogLevel,
		envGoghLogDate,
		envGoghLogTime,
		envGoghLogMicroSeconds,
		envGoghLogLongFile,
		envGoghLogShortFile,
		envGoghLogUTC,
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
	Log: LogConfig{
		Level: DefaultLogLevel,
		Time:  TrueConfig,
	},
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
		defaultConfig.VRoot = util.UniqueStringArray(root)
	})
	return &defaultConfig
}

func LoadConfig(r io.Reader) (config *Config, err error) {
	config = &Config{}
	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}
	config.VRoot = util.UniqueStringArray(config.VRoot)
	return
}

func GetEnvarConfig() (config *Config, err error) {
	config = &Config{}
	err = envdecode.Decode(config)
	if err == envdecode.ErrNoTargetFieldsAreSet {
		err = nil
	}
	config.VRoot = util.UniqueStringArray(config.VRoot)
	return
}

func MergeConfig(base *Config, override ...*Config) *Config {
	c := *base
	for _, o := range override {
		c.Log.Level = mergeStringOption(c.Log.Level, o.Log.Level)
		c.Log.Date = mergeBoolOption(c.Log.Date, o.Log.Date)
		c.Log.Time = mergeBoolOption(c.Log.Time, o.Log.Time)
		c.Log.MicroSeconds = mergeBoolOption(c.Log.MicroSeconds, o.Log.MicroSeconds)
		c.Log.LongFile = mergeBoolOption(c.Log.LongFile, o.Log.LongFile)
		c.Log.ShortFile = mergeBoolOption(c.Log.ShortFile, o.Log.ShortFile)
		c.Log.UTC = mergeBoolOption(c.Log.UTC, o.Log.UTC)
		c.VRoot = mergeStringArrayOption(c.VRoot, o.VRoot)
		c.GitHub.Token = mergeStringOption(c.GitHub.Token, o.GitHub.Token)
		c.GitHub.User = mergeStringOption(c.GitHub.User, o.GitHub.User)
		c.GitHub.Host = mergeStringOption(c.GitHub.Host, o.GitHub.Host)
	}
	return &c
}

func mergeBoolOption(base, override BoolConfig) BoolConfig {
	switch {
	case override.filled:
		return override
	case base.filled:
		return base
	default:
		return BoolConfig{}
	}
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
