package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configDir string
	cacheDir  string
	homeDir   string
)

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("search user home dir: %w", err))
	}
	configDir, err = os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("search user config dir: %w", err))
	}
	cacheDir, err = os.UserCacheDir()
	if err != nil {
		panic(fmt.Errorf("search user cache dir: %w", err))
	}
}

var Config struct {
	Roots []string
}

func GetRoots() []string {
	roots := make([]string, 0, len(Config.Roots))
	for _, r := range Config.Roots {
		r = os.ExpandEnv(r)

		runes := []rune(r)
		if runes[0] == '~' && (runes[1] == filepath.Separator || runes[1] == '/') {
			r = filepath.Join(homeDir, string(runes[2:]))
		}
		roots = append(roots, r)
	}
	return roots
}

var Servers gogh.Servers

func setupConfig(*cobra.Command, []string) error {
	var notFound viper.ConfigFileNotFoundError
	copt := viper.New()
	copt.SetConfigName("config")
	copt.AddConfigPath(filepath.Join(configDir, appname))
	copt.SetConfigType("yaml")
	copt.SetEnvPrefix(appname)
	if err := copt.ReadInConfig(); err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("read in user config: %w", err)
	}
	if err := copt.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	sopt := viper.New()
	sopt.SetConfigName("servers")
	sopt.AddConfigPath(filepath.Join(cacheDir, appname))
	sopt.SetConfigType("yaml")
	if err := sopt.ReadInConfig(); err != nil && !errors.As(err, &notFound) {
		return fmt.Errorf("read in cached servers: %#v, %w", err, err)
	}
	if err := sopt.Unmarshal(&Servers); err != nil {
		return fmt.Errorf("unmarshal cached servers: %w", err)
	}
	return nil
}
