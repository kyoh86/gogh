package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/bundle_dump"
	"github.com/kyoh86/gogh/v3/app/config"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewBundleDumpCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var flags config.BundleDumpFlags
	cmd := &cobra.Command{
		Use:     "dump",
		Aliases: []string{"export"},
		Short:   "Export current local repository list",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := os.Stdout
			if flags.File != "" {
				file, err := os.OpenFile(
					flags.File,
					os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
					0644,
				)
				if err != nil {
					return fmt.Errorf("failed to open file: %w", err)
				}
				defer file.Close()
				out = file
			}
			useCase := bundle_dump.NewUseCase(svc.WorkspaceService, svc.FinderService, svc.HostingService, svc.GitService)
			for entry, err := range useCase.Execute(cmd.Context(), bundle_dump.Options{}) {
				if err != nil {
					log.FromContext(cmd.Context()).Error(err.Error())
					return nil
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

	cmd.Flags().StringVarP(&flags.File, "file", "f", svc.Flags.BundleDump.File, "A file to output; if not specified, output to stdout")
	return cmd, nil
}
