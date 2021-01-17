package mainutil

import (
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/xdg"
)

var (
	aliasFile = filepath.Join(xdg.DataHome(), "gogh", "alias.yaml")
)

func setConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.ConfigHome(), "gogh", "config.yaml")).
		Envar("GOGH_CONFIG").
		StringVar(configFile)
}

func WrapCommand(cmd *kingpin.CmdClause, f func(gogh.Env) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() (retErr error) {
		access, err := loadAccess(configFile)
		if err != nil {
			return err
		}

		if err := loadAlias(aliasFile); err != nil {
			return err
		}

		return f(&access)
	}
}

func WrapConfigurableCommand(cmd *kingpin.CmdClause, f func(gogh.Env, *config.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() (retErr error) {
		config, access, err := loadAppenv(configFile)
		if err != nil {
			return err
		}

		if err = f(&access, &config); err != nil {
			return err
		}

		return saveConfig(configFile, config)
	}
}
