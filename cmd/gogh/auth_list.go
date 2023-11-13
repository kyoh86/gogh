package main

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var authListCommand = &cobra.Command{
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
		for _, entry := range entries {
			fmt.Println(entry)
		}
		return nil
	},
}

func init() {
	authCommand.AddCommand(authListCommand)
}
