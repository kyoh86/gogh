package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kyoh86/gogh/v2/cmdutil"
	"github.com/spf13/cobra"
)

var setupOnce sync.Once

func setup() {
	setupOnce.Do(func() {
		cobra.CheckErr(setupCore())
	})
}

var (
	configFilePath      string
	tokensFilePath      string
	defaultFlagFilePath string
)

func setupCore() (err error) {
	configFilePath, err = loadConfig()
	if err != nil {
		return
	}
	tokensFilePath, err = loadTokens()
	if err != nil {
		return
	}
	defaultFlagFilePath, err = loadDefaultFlag()
	if err != nil {
		return
	}
	return nil
}

func envOrDir(envName string, dirHandler func() (string, error)) func() (string, error) {
	return func() (string, error) {
		if env := os.Getenv(envName); env != "" {
			return env, nil
		}
		return dirHandler()
	}
}

var (
	configFileHandler      = cmdutil.AppFile{EnvName: "GOGH_CONFIG_PATH", Dir: os.UserConfigDir, Basename: "config.yaml"}
	defaultFlagFileHandler = cmdutil.AppFile{EnvName: "GOGH_FLAG_PATH", Dir: os.UserConfigDir, Basename: "flag.yaml"}
	tokensFileHandler      = cmdutil.AppFile{EnvName: "GOGH_TOKENS_PATH", Dir: os.UserCacheDir, Basename: "tokens.yaml"}
)

func loadConfig() (string, error) {
	path, err := configFileHandler.Load(&config)
	if err != nil {
		return "", err
	}
	if len(config.Roots) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("search user home dir: %w", err)
		}
		raw := filepath.Join(homeDir, "Projects")
		config.Roots = []expandedPath{{
			raw:      raw,
			expanded: raw,
		}}
	}
	return path, nil
}

func saveConfig() error {
	return configFileHandler.Save(config)
}

func loadTokens() (string, error) {
	return tokensFileHandler.Load(&tokens)
}

func saveTokens() error {
	return tokensFileHandler.Save(tokens)
}

func loadDefaultFlag() (string, error) {
	return defaultFlagFileHandler.Load(&defaultFlag)
}
