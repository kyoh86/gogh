package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewExtraCommand(_ context.Context, _ *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "extra",
		Short: "Manage repository extra files",
		Long: `Manage extra files that are typically ignored by git.

There are two types of extra:
1. Auto-apply extra: Automatically applied when cloning the repository
2. Named extra: Templates that can be manually applied to any repository`,
	}
	return cmd, nil
}
