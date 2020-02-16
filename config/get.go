package config

import (
	"go/build"
	"io"
	"log"
	"path/filepath"
	"sync"

	"github.com/joeshaw/envdecode"
	"github.com/thoas/go-funk"
	"github.com/zalando/go-keyring"
	yaml "gopkg.in/yaml.v3"
)

var (
	envGoghGitHubUser  = "GOGH_GITHUB_USER"
	envGoghGitHubToken = "GOGH_GITHUB_TOKEN"
	envGoghGitHubHost  = "GOGH_GITHUB_HOST"
	envGoghRoot        = "GOGH_ROOT"
	envNames           = []string{
		envGoghGitHubUser,
		envGoghGitHubToken,
		envGoghGitHubHost,
		envGoghRoot,
	}
	keyGoghServiceName = "gogh.kyoh86.dev"
	keyGoghGitHubToken = "github-token"
)

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost = "github.com"
)

var defaultConfig = Config{
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

func LoadKeyring() *Config {
	token, err := keyring.Get(keyGoghServiceName, keyGoghGitHubToken)
	if err != nil {
		log.Printf("info: there's no token in %s::%s (%v)", keyGoghServiceName, keyGoghGitHubToken, err)
		return &Config{}
	}

	return &Config{GitHub: GitHubConfig{Token: token}}
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
