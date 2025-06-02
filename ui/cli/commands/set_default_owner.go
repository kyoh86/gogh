package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewSetDefaultOwnerCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "set-default-owner [flags] <host> <owner>",
		Short: "Set the default owner for a host for the repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			host, owner := args[0], args[1]
			if err := svc.DefaultNameService.SetDefaultOwnerFor(host, owner); err != nil {
				return fmt.Errorf("setting default host: %w", err)
			}
			log.FromContext(ctx).Infof("Default owner for %s host set to %s\n", host, owner)
			return nil
		},
	}
	return cmd, nil
}
