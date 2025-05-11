package commands

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/cwd"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewCwdCommand(svc *ServiceSet) *cobra.Command {
	var f config.CwdFlags
	cmd := &cobra.Command{
		Use:   "cwd",
		Short: "Print the local reposiotry in current working directory",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := f.Format.Formatter()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			repo, err := cwd.NewUseCase(svc.workspaceService, svc.finderService).Execute(ctx, wd)
			if err != nil {
				return err
			}
			str, err := formatter.Format(repo)
			if err != nil {
				log.FromContext(ctx).WithFields(log.Fields{
					"error":  err,
					"format": f.Format.String(),
					"path":   repo.FullPath(),
				}).Info("failed to format")
			}
			fmt.Println(str)
			return nil
		},
	}

	f.Format = svc.defaults.Cwd.Format
	cmd.Flags().VarP(&f.Format, "format", "f", flags.LocalRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteLocalRepoFormat); err != nil {
		panic(err)
	}
	return cmd
}
