package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewSetDefaultOwnerCommand(svc *ServiceSet) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-default-owner <host> <owner>",
		Short: "Set the default owner for a host for the repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			host, owner := args[0], args[1]
			if err := svc.defaultNameService.SetDefaultOwnerFor(host, owner); err != nil {
				return fmt.Errorf("failed to set default host: %w", err)
			}
			fmt.Printf("Default host set to %s\n", host)
			return nil
		},
	}
	return cmd
}
