package mainutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/xdg"
)

func setConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.ConfigHome(), "gogh", "config.yaml")).
		Envar("GOGH_CONFIG").
		StringVar(configFile)
}

func openYAML(filename string) (io.Reader, func() error, error) {
	var reader io.Reader
	var teardown func() error
	file, err := os.Open(filename)
	switch {
	case err == nil:
		teardown = func() error { return file.Close() }
		reader = file
	case os.IsNotExist(err):
		reader = env.EmptyYAMLReader
		teardown = func() error { return nil }
	default:
		return nil, nil, err
	}
	return reader, teardown, nil
}

func WrapCommand(cmd *kingpin.CmdClause, f func(gogh.Env) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		reader, teardown, err := openYAML(configFile)
		if err != nil {
			return err
		}
		defer teardown()
		access, err := env.GetAccess(reader, env.KeyringService, env.EnvarPrefix)
		if err != nil {
			return err
		}

		return f(&access)
	}
}

func WrapConfigurableCommand(cmd *kingpin.CmdClause, f func(*env.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		reader, teardown, err := openYAML(configFile)
		if err != nil {
			return err
		}
		defer teardown()
		config, err := env.GetConfig(reader, env.KeyringService)
		if err != nil {
			return err
		}

		if err = f(&config); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(configFile), 0744); err != nil {
			return err
		}
		file, err := os.OpenFile(configFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		return config.Save(file, env.KeyringService)
	}
}
