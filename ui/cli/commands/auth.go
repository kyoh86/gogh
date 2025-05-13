package commands

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewAuthCommand(_ context.Context, svc *service.ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
	}
}
