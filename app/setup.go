package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	Name = "gogh"
)

var (
	homeDir   string
	configDir string
	cacheDir  string
)

func Setup() error {
	{
		dir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("search user home dir: %w", err)
		}
		homeDir = dir
	}
	{
		dir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("search user config dir: %w", err)
		}
		configDir = dir
	}
	{
		dir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("search user cache dir: %w", err)
		}
		cacheDir = dir
	}

	if err := setupConfig(); err != nil {
		return err
	}

	if err := setupServers(); err != nil {
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
