package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/level"
	"github.com/apex/log/handlers/multi"
	"github.com/kyoh86/gogh/v3/config"
	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

// StdoutLogHandler implementation.
type StdoutLogHandler struct {
	Handler log.Handler
}

// HandleLog implements log.Handler.
func (h *StdoutLogHandler) HandleLog(e *log.Entry) error {
	if e.Level >= log.ErrorLevel {
		return nil
	}

	return h.Handler.HandleLog(e)
}

func main() {
	errLog := level.New(cli.New(os.Stderr), log.ErrorLevel)
	stdLog := &StdoutLogHandler{Handler: cli.New(os.Stdout)}
	level := log.InfoLevel
	if os.Getenv("GOGH_DEBUG") == "1" {
		level = log.DebugLevel
	}
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: multi.New(stdLog, errLog),
		Level:   level,
	})

	flags, err := config.LoadFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load flags: %s\n", err)
		os.Exit(1)
	}
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s\n", err)
		os.Exit(1)
	}
	tokens, err := config.LoadTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load tokens: %s\n", err)
		os.Exit(1)
	}

	facadeCommand := &cobra.Command{
		Use:     config.AppName,
		Short:   "GO GitHub project manager",
		Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
	}

	bundleCommand := NewBundleCommand()
	bundleCommand.AddCommand(
		NewBundleDumpCommand(conf, flags),
		NewBundleRestoreCommand(conf, tokens, flags),
	)

	authCommand := NewAuthCommand()
	authCommand.AddCommand(
		NewAuthListCommand(tokens),
		NewAuthLoginCommand(tokens),
		NewAuthLogoutCommand(tokens),
		NewAuthSetDefaultCommand(tokens),
	)

	rootsCommand := NewRootsCommand(conf)
	rootsCommand.AddCommand(
		NewRootsSetDefaultCommand(conf),
		NewRootsRemoveCommand(conf),
		NewRootsAddCommand(conf),
		NewRootsListCommand(conf),
	)

	configCommand := NewConfigCommand(conf, tokens, flags)
	configCommand.AddCommand(
		authCommand,
		rootsCommand,
	)

	facadeCommand.AddCommand(
		NewCwdCommand(conf, flags),
		NewListCommand(conf, flags),
		NewCloneCommand(conf, tokens),
		NewCreateCommand(conf, tokens, flags),
		NewReposCommand(tokens, flags),
		NewDeleteCommand(conf, tokens),
		NewForkCommand(conf, tokens, flags),
		configCommand,
		authCommand,
		bundleCommand,
		rootsCommand,
	)

	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		log.FromContext(ctx).Error(err.Error())
		os.Exit(1)
	}
}
