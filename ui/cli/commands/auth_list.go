package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/auth/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewAuthListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "list",
		Short: "Listup authenticated host and owners",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			entries, err := list.NewUsecase(svc.TokenService).Execute(ctx)
			if err != nil {
				return fmt.Errorf("listing up tokens: %w", err)
			}
			if len(entries) == 0 {
				return errors.New("no valid token found: you need to set token by `gogh auth login`")
			}
			for _, entry := range entries {
				fmt.Printf("%s/%s\n", entry.Host, entry.Owner)
			}
			return nil
		},
	}, nil
}
