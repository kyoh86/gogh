package main

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"
)

type bundleRestoreFlagsStruct struct {
	File   expandedPath `yaml:"file,omitempty"`
	Dryrun bool         `yaml:"-"`
}

var (
	bundleRestoreFlags   bundleRestoreFlagsStruct
	bundleRestoreCommand = &cobra.Command{
		Use:   "restore",
		Short: "Get dumped projects",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			in := os.Stdin
			if bundleRestoreFlags.File.expanded != "" {
				f, err := os.Open(bundleRestoreFlags.File.expanded)
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
			return cloneAll(ctx, specs, bundleRestoreFlags.Dryrun)
		},
	}
)

func init() {
	setup()
	bundleRestoreFlags.File = defaultFlag.BundleRestore.File
	bundleRestoreCommand.Flags().BoolVarP(&bundleRestoreFlags.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	bundleRestoreCommand.Flags().VarP(&bundleRestoreFlags.File, "file", "", "Read the file as input")
	bundleCommand.AddCommand(bundleRestoreCommand)
}
