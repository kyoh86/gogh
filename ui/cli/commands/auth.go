package commands

import (
	"context"

	"github.com/spf13/cobra"
)

func NewAuthCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
	}
}
