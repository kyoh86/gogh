package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/bundle_dump"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewBundleDumpCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.BundleDumpFlags
	cmd := &cobra.Command{
		Use:     "dump",
		Aliases: []string{"export"},
		Short:   "Export current local repository list",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := cmd.OutOrStdout()
			if f.File != "" && f.File != "-" {
				file, err := os.OpenFile(
					f.File,
					os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
					0644,
				)
				if err != nil {
					return fmt.Errorf("opening file: %w", err)
				}
				defer file.Close()
				out = file
			}
			for entry, err := range bundle_dump.NewUsecase(svc.WorkspaceService, svc.FinderService, svc.HostingService, svc.GitService).Execute(cmd.Context(), bundle_dump.Options{}) {
				if err != nil {
					return err
				}
				if entry.Alias == nil {
					fmt.Fprintln(out, entry.Name)
				} else {
					fmt.Fprintf(out, "%s=%s\n", *entry.Alias, entry.Name)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&f.File, "file", "f", svc.Flags.BundleDump.File, `A file to output; if it's empty("") or hyphen("-"), output to stdout`)
	return cmd, nil
}
