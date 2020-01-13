package mainutil

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/comail/colog"
	"github.com/kyoh86/gogh/config"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/xdg"
)

func setConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.ConfigHome(), "gogh", "config.yaml")).
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
	colog.SetFlags(ctx.LogFlags())
	colog.SetMinLevel(lvl)
	colog.SetDefaultLevel(colog.LError)
	return nil
}

func currentConfig(configFile string, validate bool) (*config.Config, *config.Config, error) {
	var savedConfig *config.Config
	file, err := os.Open(configFile)
	switch {
	case err == nil:
		defer file.Close()
		savedConfig, err = config.LoadConfig(file)
		if err != nil {
			return nil, nil, err
		}
	case os.IsNotExist(err):
		savedConfig = &config.Config{}
	default:
		return nil, nil, err
	}

	savedConfig = config.MergeConfig(savedConfig, config.LoadKeyring())
	envarConfig, err := config.GetEnvarConfig()
	if err != nil {
		return nil, nil, err
	}
	cfg := config.MergeConfig(config.DefaultConfig(), savedConfig, envarConfig)
	if !validate {
		return savedConfig, cfg, nil
	}
	if err := gogh.ValidateContext(cfg); err != nil {
		return nil, nil, err
	}
	return savedConfig, cfg, nil
}

func WrapCommand(cmd *kingpin.CmdClause, f func(gogh.Context) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		_, cfg, err := currentConfig(configFile, true)
		if err != nil {
			return err
		}

		if err := initLog(cfg); err != nil {
			return err
		}
		return f(cfg)
	}
}

func WrapConfigurableCommand(cmd *kingpin.CmdClause, f func(*config.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		savedConfig, cfg, err := currentConfig(configFile, false)
		if err != nil {
			return err
		}

		if err := initLog(cfg); err != nil {
			return err
		}

		if err = f(savedConfig); err != nil {
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
		return config.SaveConfig(file, savedConfig)
	}
}
