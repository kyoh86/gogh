package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/auth_list"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewAuthListCommand(_ context.Context, svc *service.ServiceSet) *cobra.Command {
	useCase := auth_list.NewUseCase(svc.TokenService)
	return &cobra.Command{
		Use:   "list",
		Short: "Listup authenticated host and owners",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			entries, err := useCase.Execute(ctx)
			if err != nil {
				log.FromContext(ctx).WithError(err).Error("failed to list tokens")
				return nil
			}
			if len(entries) == 0 {
				log.FromContext(ctx).Warn("No valid token found: you need to set token by `gogh auth login`")
				return nil
			}
			for _, entry := range entries {
				fmt.Printf("  %s%s\n", entry.Host, entry.Owner)
			}
			return nil
		},
	}
}
