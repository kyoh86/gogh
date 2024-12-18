package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/level"
	"github.com/apex/log/handlers/multi"
	"github.com/kyoh86/gogh/v2/cmdutil"
	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

var facadeCommand = &cobra.Command{
	Use:     cmdutil.AppName,
	Short:   "GO GitHub project manager",
	Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
}

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
	setup()
	errLog := level.New(cli.New(os.Stderr), log.ErrorLevel)
	stdLog := &StdoutLogHandler{Handler: cli.New(os.Stderr)}
	level := log.InfoLevel
	if os.Getenv("GOGH_DEBUG") == "1" {
		level = log.DebugLevel
	}
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: multi.New(stdLog, errLog),
		Level:   level,
	})
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		log.FromContext(ctx).Error(err.Error())
		os.Exit(1)
	}
}
