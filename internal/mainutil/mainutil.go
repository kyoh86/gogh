package mainutil

import (
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/comail/colog"
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/xdg"
)

func SetConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.CacheHome(), "gogh", "config.toml")).
		Envar("GOGH_CONFIG").
		StringVar(configFile)
}

func InitLog(ctx gogh.Context) error {
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

func WrapCommand(cmd *kingpin.CmdClause, f func(gogh.Context) error) (string, func() error) {
	var configFile string
	SetConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		fileConfig, err := command.LoadFileConfig(configFile)
		if err != nil {
			return err
		}
		envarConfig, err := command.GetEnvarConfig()
		if err != nil {
			return err
		}

		ctx := command.MergeConfig(command.DefaultConfig(), fileConfig, envarConfig)

		InitLog(&ctx)
		return f(&ctx)
	}
}
