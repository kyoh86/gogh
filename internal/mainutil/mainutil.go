package mainutil

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
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

func currentConfig(configFile string) (*config.Config, *config.Config, error) {
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
	if err := gogh.ValidateContext(cfg); err != nil {
		log.Printf("warn: invalid config: %v", err)
	}
	return savedConfig, cfg, nil
}

func WrapCommand(cmd *kingpin.CmdClause, f func(gogh.Context) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		_, cfg, err := currentConfig(configFile)
		if err != nil {
			return err
		}

		return f(cfg)
	}
}

func WrapConfigurableCommand(cmd *kingpin.CmdClause, f func(*config.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		savedConfig, _, err := currentConfig(configFile)
		if err != nil {
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
