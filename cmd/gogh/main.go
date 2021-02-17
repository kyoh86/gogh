package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

const appname = "gogh"

var facadeCommand = &cobra.Command{
	Use:               appname,
	Short:             "GO GitHub project manager",
	Version:           fmt.Sprintf("%s-%s (%s)", version, commit, date),
	PersistentPreRunE: setupConfig,
}

func main() {
	if err := facadeCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
