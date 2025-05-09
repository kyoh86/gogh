package cmdutil

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

type AppFile struct {
	EnvName  string
	Dir      func() (string, error)
	Basename string
}

func (h AppFile) appFilePath() (string, error) {
	if env := os.Getenv(h.EnvName); env != "" {
		return env, nil
	}
	dir, err := h.Dir()
	if err != nil {
		return "", fmt.Errorf("search app file dir for %s: %w", h.Basename, err)
	}
	return filepath.Join(dir, AppName, h.Basename), nil
}

func (h AppFile) Load(output any) (string, error) {
	path, err := h.appFilePath()
	if err != nil {
		return "", err
	}
	if err := loadYAML(path, output); err != nil {
		return "", fmt.Errorf("load %s: %w", h.Basename, err)
	}
	return path, nil
}

func (h AppFile) Save(input any) error {
	path, err := h.appFilePath()
	if err != nil {
		return err
	}
	if err := saveYAML(path, input); err != nil {
		return fmt.Errorf("save %s: %w", h.Basename, err)
	}
	return nil
}
