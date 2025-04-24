package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "github.com/goccy/go-yaml"
)

const AppName = "gogh"

func loadYAML(path string, obj any) error {
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

func saveYAML(path string, obj any) error {
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

func appFilePath(envar string, getDir func() (string, error), rel ...string) (string, error) {
	if env := os.Getenv(envar); env != "" {
		return env, nil
	}
	dir, err := getDir()
	if err != nil {
		return "", fmt.Errorf("search app file dir for %s: %w", rel, err)
	}
	return filepath.Join(append([]string{dir, AppName}, rel...)...), nil
}
