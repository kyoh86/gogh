// Package cmd includes processes of all subcommands
package cmd

import "context"

// Command will runs when the flags matched
type Command func(ctx context.Context) error
