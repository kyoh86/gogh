package commands

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewSetDefaultHostCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "set-default-host",
		Short: "Set the default host for the repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			if err := svc.DefaultNameService.SetDefaultHost(host); err != nil {
				return fmt.Errorf("setting default host: %w", err)
			}
			fmt.Printf("Default host set to %s\n", host)
			return nil
		},
	}
	return cmd, nil
}
