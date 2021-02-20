package main

import (
	"fmt"
	"os"

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
	if err := facadeCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
