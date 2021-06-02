package main

import (
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
