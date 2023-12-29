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

var (
	configFileHandler      = cmdutil.AppFile{Dir: os.UserConfigDir, Basename: "config.yaml"}
	tokensFileHandler      = cmdutil.AppFile{Dir: os.UserCacheDir, Basename: "tokens.yaml"}
	defaultFlagFileHandler = cmdutil.AppFile{Dir: os.UserConfigDir, Basename: "flag.yaml"}
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
	if tokens.DefaultHost == "" {
		tokens.DefaultHost = "github.com"
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
