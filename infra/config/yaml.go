package config

import (
	"os"
	"path/filepath"

	yaml "github.com/goccy/go-yaml"
)

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
