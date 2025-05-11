package commands

import (
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/auth_list"
	"github.com/spf13/cobra"
)

func NewAuthListCommand(svc *ServiceSet) *cobra.Command {
	useCase := auth_list.NewUseCase(svc.tokenService)
	return &cobra.Command{
		Use:   "list",
		Short: "Listup authenticated host and owners",
		Args:  cobra.ExactArgs(0),
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
