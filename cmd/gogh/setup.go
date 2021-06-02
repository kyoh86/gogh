package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

const appName = "gogh"

var setupOnce sync.Once

func setup() {
	setupOnce.Do(func() {
		cobra.CheckErr(setupCore())
	})
}

func setupCore() error {
	if err := loadConfig(); err != nil {
		return err
	}

	if err := loadServers(); err != nil {
		return err
	}

	if err := loadDefaultFlag(); err != nil {
		return err
	}

	return nil
}

func loadYAML(path string, obj interface{}) error {
	file, err := os.Open(path)
	switch {
	case err == nil:
		// noop
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
	defer file.Close()
	dec := yaml.NewDecoder(file)
	return dec.Decode(obj)
}

func saveYAML(path string, obj interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := yaml.NewEncoder(file)
	return enc.Encode(obj)
}

type appFileHandler struct {
	dir      func() (string, error)
	basename string
}

func (h appFileHandler) load(output interface{}) error {
	dir, err := h.dir()
	if err != nil {
		return fmt.Errorf("search app file dir for %s: %w", h.basename, err)
	}
	path := filepath.Join(dir, appName, h.basename)
	if err := loadYAML(path, output); err != nil {
		return fmt.Errorf("load %s: %w", h.basename, err)
	}
	return nil
}

func (h appFileHandler) save(input interface{}) error {
	dir, err := h.dir()
	if err != nil {
		return fmt.Errorf("search app file dir for %s: %w", h.basename, err)
	}
	path := filepath.Join(dir, appName, h.basename)
	if err := saveYAML(path, input); err != nil {
		return fmt.Errorf("save %s: %w", h.basename, err)
	}
	return nil
}

var (
	configFileHandler      = appFileHandler{dir: os.UserConfigDir, basename: "config.yaml"}
	serversFileHandler     = appFileHandler{dir: os.UserCacheDir, basename: "servers.yaml"}
	defaultFlagFileHandler = appFileHandler{dir: os.UserConfigDir, basename: "flag.yaml"}
)

func loadConfig() error {
	if err := configFileHandler.load(&config); err != nil {
		return err
	}
	if len(config.Roots) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("search user home dir: %w", err)
		}
		raw := filepath.Join(homeDir, "Projects")
		config.Roots = []expandedPath{{
			raw:      raw,
			expanded: raw,
		}}
	}
	return nil
}

func saveConfig() error {
	return configFileHandler.save(config)
}

func loadServers() error {
	return serversFileHandler.load(&servers)
}

func saveServers() error {
	return serversFileHandler.save(servers)
}

func loadDefaultFlag() error {
	return defaultFlagFileHandler.load(&defaultFlag)
}
