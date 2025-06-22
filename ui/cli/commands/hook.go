package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewHookCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Manage repository hooks",
	}
	return cmd, nil
}
