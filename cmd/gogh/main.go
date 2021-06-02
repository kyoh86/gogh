package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

var facadeCommand = &cobra.Command{
	Use:     Name,
	Short:   "GO GitHub project manager",
	Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
}

func main() {
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	})
	if err := Setup(); err != nil {
		log.FromContext(ctx).WithField("error", err).Fatalf("failed to setup application")
	}
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
