package mainutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh/alias"
)

func openYAML(filename string) (io.Reader, func() error, error) {
	var reader io.Reader
	var teardown func() error
	file, err := os.Open(filename)
	switch {
	case err == nil:
		teardown = file.Close
		reader = file
	case os.IsNotExist(err):
		reader = config.EmptyYAMLReader
		teardown = func() error { return nil }
	default:
		return nil, nil, err
	}
	return reader, teardown, nil
}

func loadAccess(configFile string) (_ config.Access, retErr error) {
	reader, teardown, err := openYAML(configFile)
	if err != nil {
		retErr = err
		return
	}
	defer func() {
		if err := teardown(); err != nil && retErr == nil {
			retErr = err
			return
		}
	}()
	return config.GetAccess(reader, config.EnvarPrefix)
}

func loadAppenv(configFile string) (_ config.Config, _ config.Access, retErr error) {
	reader, teardown, err := openYAML(configFile)
	if err != nil {
		retErr = err
		return
	}
	defer func() {
		if err := teardown(); err != nil && retErr == nil {
			retErr = err
			return
		}
	}()
	return config.GetAppenv(reader, config.EnvarPrefix)
}

func saveConfig(configFile string, c config.Config) error {
	if err := os.MkdirAll(filepath.Dir(configFile), 0744); err != nil {
		return err
	}
	file, err := os.OpenFile(configFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return c.Save(file)
}

func loadAlias(configFile string) (retErr error) {
	reader, teardown, err := openYAML(configFile)
	if err != nil {
		retErr = err
		return
	}
	defer func() {
		if err := teardown(); err != nil && retErr == nil {
			retErr = err
			return
		}
	}()
	var d alias.Def
	if err := yaml.NewDecoder(reader).Decode(&d); err != nil {
		return err
	}
	alias.Instance = d // nolint
	return nil
}
