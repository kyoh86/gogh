package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewAuthListCommand(tokens *config.TokenStore) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Listup authenticated host and owners",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			entries := tokens.Entries()
			if len(entries) == 0 {
				log.FromContext(ctx).Warn("No valid token found: you need to set token by `gogh auth login`")
				return nil
			}
			host, owner := tokens.GetDefaultKey()
			for _, entry := range entries {
				if entry.Host == host && entry.Owner == owner {
					fmt.Printf("* %s\n", entry)
				} else {
					fmt.Printf("  %s\n", entry)
				}
			}
			return nil
		},
	}
}
