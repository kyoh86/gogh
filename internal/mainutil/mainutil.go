package mainutil

import (
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/comail/colog"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/xdg"
)

func setConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.CacheHome(), "gogh", "config.toml")).
		Envar("GOGH_CONFIG").
		StringVar(configFile)
}

func initLog(ctx gogh.Context) error {
	lvl, err := colog.ParseLevel(ctx.LogLevel())
	if err != nil {
		return err
	}
	colog.Register()
	colog.SetOutput(ctx.Stderr())
	colog.SetMinLevel(lvl)
	colog.SetDefaultLevel(colog.LError)
	return nil
}

func currentConfig(configFile string) (*gogh.Config, *gogh.Config, error) {
	fileConfig, err := gogh.LoadFileConfig(configFile)
	if err != nil {
		return nil, nil, err
	}
	envarConfig, err := gogh.GetEnvarConfig()
	if err != nil {
		return nil, nil, err
	}
	config := gogh.MergeConfig(gogh.DefaultConfig(), fileConfig, envarConfig)
	if err := gogh.ValidateContext(config); err != nil {
		return nil, nil, err
	}
	return fileConfig, config, nil
}

func WrapCommand(cmd *kingpin.CmdClause, f func(gogh.Context) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		_, config, err := currentConfig(configFile)
		if err != nil {
			return err
		}

		if err := initLog(config); err != nil {
			return err
		}
		return f(config)
	}
}

func WrapConfigurableCommand(cmd *kingpin.CmdClause, f func(gogh.Context, *gogh.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		fileConfig, config, err := currentConfig(configFile)
		if err != nil {
			return err
		}

		if err := initLog(config); err != nil {
			return err
		}

		return f(config, fileConfig)
	}
}
