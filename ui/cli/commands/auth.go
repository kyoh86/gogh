package commands

import (
	"github.com/spf13/cobra"
)

func NewAuthCommand(svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
	}
}
