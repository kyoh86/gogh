package commands

import (
	"context"
	_ "embed"

	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewConfigCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "config",
		Short:   "Show configurations",
		Aliases: []string{"conf", "setting", "context"},
	}, nil
}
