package commands

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/bundle_dump"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewBundleDumpCommand(svc *ServiceSet) *cobra.Command {
	var f config.BundleDumpFlags
	cmd := &cobra.Command{
		Use:     "dump",
		Aliases: []string{"export"},
		Short:   "Export current local repository list",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := os.Stdout
			if f.File.Expand() != "" {
				f, err := os.OpenFile(
					f.File.Expand(),
					os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
					0644,
				)
				if err != nil {
					return err
				}
				defer f.Close()
				out = f
			}
			useCase := bundle_dump.NewUseCase(svc.workspaceService, svc.finderService, svc.gitService)
			for entry, err := range useCase.Execute(cmd.Context(), workspace.ListOptions{}) {
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

	f.File = svc.defaults.BundleDump.File
	cmd.Flags().VarP(&f.File, "file", "", "A file to output; if not specified, output to stdout")
	return cmd
}
