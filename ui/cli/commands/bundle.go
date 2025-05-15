package commands

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewBundleCommand(_ context.Context, _ *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "bundle",
		Short: "Manage bundle",
	}, nil
}
