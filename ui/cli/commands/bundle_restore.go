package commands

import (
	"bufio"
	"os"

	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewBundleRestoreCommand(conf *config.ConfigStore, tokens *config.TokenStore, defaults *config.FlagStore) *cobra.Command {
	var f config.BundleRestoreFlags
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Get dumped projects",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			in := os.Stdin
			if f.File.Expand() != "" {
				f, err := os.Open(f.File.Expand())
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
			return cloneAll(ctx, conf, tokens, specs, f.Dryrun)
		},
	}
	cmd.Flags().
		BoolVarP(&f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	f.File = defaults.BundleRestore.File
	cmd.Flags().
		VarP(&f.File, "file", "", "Read the file as input; if not specified, read from stdin")
	return cmd
}
