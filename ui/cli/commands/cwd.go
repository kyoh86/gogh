package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/cwd"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewCwdCommand(ctx context.Context, svc *ServiceSet) *cobra.Command {
	var format flags.LocationFormat

	cmd := &cobra.Command{
		Use:   "cwd",
		Short: "Print the local reposiotry in current working directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := format.Formatter()
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
			str, err := formatter.Format(*repo)
			if err != nil {
				log.FromContext(ctx).WithFields(log.Fields{
					"error":  err,
					"format": format.String(),
					"path":   repo.FullPath(),
				}).Info("failed to format")
			}
			fmt.Println(str)
			return nil
		},
	}

	if err := flags.LocationFormatFlag(cmd, &format, svc.flags.Cwd.Format); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to init format flag")
	}
	return cmd
}
