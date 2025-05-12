package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSetDefaultHostCommand(svc *ServiceSet) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-default-host",
		Short: "Set the default host for the repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			if err := svc.defaultNameService.SetDefaultHost(host); err != nil {
				return fmt.Errorf("failed to set default host: %w", err)
			}
			fmt.Printf("Default host set to %s\n", host)
			return nil
		},
	}
	return cmd
}
