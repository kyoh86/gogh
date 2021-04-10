package main

import (
	"bufio"
	"os"

	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var bundleRestoreFlags struct {
	file   string
	dryrun bool
}

var bundleRestoreCommand = &cobra.Command{
	Use:   "restore",
	Short: "Get dumped projects",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		in := os.Stdin
		if bundleRestoreFlags.file != "" {
			f, err := os.Open(bundleRestoreFlags.file)
			if err != nil {
				return err
			}
			defer f.Close()
			in = f
		}
		var specs []string
		scan := bufio.NewScanner(in)
		for scan.Scan() {
			specs = append(specs, scan.Text())
		}

		ctx := cmd.Context()
		servers := app.Servers()
		return cloneAll(ctx, servers, specs, bundleRestoreFlags.dryrun)
	},
}

func init() {
	bundleRestoreCommand.Flags().BoolVarP(&bundleRestoreFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	bundleRestoreCommand.Flags().StringVarP(&bundleRestoreFlags.file, "file", "", "", "Read the file as input")
	bundleCommand.AddCommand(bundleRestoreCommand)
}
