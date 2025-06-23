package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
)

// NewCwdCommand creates a new command to print the local repository which the current working directory belongs to.
func NewCwdCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var format flags.LocationFormat

	cmd := &cobra.Command{
		Use:   "cwd",
		Short: "Print the local repository which the current working directory belongs to",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			formatter, err := config.LocationFormatter(format.String())
			if err != nil {
				return fmt.Errorf("invalid format: %w", err)
			}

			ctx := cmd.Context()
			repo, err := cwd.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
			if err != nil {
				return fmt.Errorf("finding repository in current directory: %w", err)
			}
			str, err := formatter.Format(*repo)
			if err != nil {
				log.FromContext(ctx).WithFields(log.Fields{
					"error":  err,
					"format": format.String(),
					"path":   repo.FullPath(),
				}).Info("Failed to format")
				return nil
			}
			fmt.Println(str)
			return nil
		},
	}

	if err := flags.LocationFormatFlag(cmd, &format, svc.Flags.Cwd.Format); err != nil {
		return nil, fmt.Errorf("adding location format flag: %w", err)
	}
	return cmd, nil
}
