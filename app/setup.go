package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/viper"
)

const (
	Name = "gogh"
)

var homeDir string
var config struct {
	Roots []string
}
var servers gogh.Servers

func Servers() *gogh.Servers {
	return &servers
}

func Roots() []string {
	roots := make([]string, 0, len(config.Roots))
	for _, r := range config.Roots {
		r = os.ExpandEnv(r)

		runes := []rune(r)
		if runes[0] == '~' && (runes[1] == filepath.Separator || runes[1] == '/') {
			r = filepath.Join(homeDir, string(runes[2:]))
		}
		roots = append(roots, r)
	}
	return roots
}

func Setup() error {
	{
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("search user home dir: %w", err)
		}
		homeDir = dir
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("search user config dir: %w", err)
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("search user cache dir: %w", err)
	}

	var notFound viper.ConfigFileNotFoundError
	copt := viper.New()
	copt.SetConfigName("config")
	copt.AddConfigPath(filepath.Join(configDir, Name))
	copt.SetConfigType("yaml")
	copt.SetEnvPrefix(Name)
	if err := copt.ReadInConfig(); err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("read in user config: %w", err)
	}
	if err := copt.Unmarshal(&config); err != nil {
		return fmt.Errorf("unmarshal cached servers: %w", err)
	}

	sopt := viper.New()
	sopt.SetConfigName("servers")
	sopt.AddConfigPath(filepath.Join(cacheDir, Name))
	sopt.SetConfigType("yaml")
	if err := sopt.ReadInConfig(); err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("read in cached servers: %#v, %w", err, err)
	}
	if err := sopt.Unmarshal(&servers); err != nil {
		return fmt.Errorf("unmarshal cached servers: %w", err)
	}
	return nil
}
