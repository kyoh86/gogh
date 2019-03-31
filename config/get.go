package config

import (
	"go/build"
	"io"
	"path/filepath"
	"sync"

	"github.com/joeshaw/envdecode"
	"github.com/thoas/go-funk"
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
		Time:  TrueOption,
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
		defaultConfig.VRoot = funk.UniqString(root)
	})
	return &defaultConfig
}

func LoadConfig(r io.Reader) (config *Config, err error) {
	config = &Config{}
	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}
	config.VRoot = funk.UniqString(config.VRoot)
	return
}

func SaveConfig(w io.Writer, config *Config) error {
	return yaml.NewEncoder(w).Encode(config)
}

func GetEnvarConfig() (config *Config, err error) {
	config = &Config{}
	err = envdecode.Decode(config)
	if err == envdecode.ErrNoTargetFieldsAreSet {
		err = nil
	}
	config.VRoot = funk.UniqString(config.VRoot)
	return
}
