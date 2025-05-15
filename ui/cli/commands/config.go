package commands

import (
	"context"
	_ "embed"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewConfigCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "config",
		Short:   "Show/change configurations",
		Aliases: []string{"conf", "setting", "context"},
	}, nil
}
