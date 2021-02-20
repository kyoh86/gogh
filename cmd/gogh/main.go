package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

var facadeCommand = &cobra.Command{
	Use:     app.Name,
	Short:   "GO GitHub project manager",
	Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
	PersistentPreRunE: func(*cobra.Command, []string) error {
		return app.Setup()
	},
}

func main() {
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	})
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
