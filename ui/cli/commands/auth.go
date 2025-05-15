package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewAuthCommand(_ context.Context, _ *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "auth",
		Short: "Manage tokens",
	}, nil
}
