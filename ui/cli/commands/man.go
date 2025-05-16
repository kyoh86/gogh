package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewManCommand(_ context.Context, _ *service.ServiceSet) (*cobra.Command, error) {
	var (
		manPagePath  string
		usageDocPath string
	)
	cmd := &cobra.Command{
		Use:    "man",
		Short:  "Generate man pages and markdown usages",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			facadeCommand := cmd.Parent()
			list, _, err := facadeCommand.Traverse([]string{"list"})
			if err != nil {
				return fmt.Errorf("finding 'list' command': %w", err)
			}
			cwd, _, err := facadeCommand.Traverse([]string{"cwd"})
			if err != nil {
				return fmt.Errorf("finding 'cwd' command: %w", err)
			}
			list.Flag("format").Usage = flags.LocationFormatLongUsage
			cwd.Flag("format").Usage = flags.LocationFormatLongUsage
			header := &doc.GenManHeader{
				Title:   "GOGH",
				Section: "1",
			}
			if err := os.MkdirAll(manPagePath, 0755); err != nil {
				return fmt.Errorf("making man page directory: %w", err)
			}
			if err := doc.GenManTree(facadeCommand, header, manPagePath); err != nil {
				return fmt.Errorf("generating man pages: %w", err)
			}
			if err := os.MkdirAll(usageDocPath, 0755); err != nil {
				return fmt.Errorf("making usage doc directory: %w", err)
			}
			facadeCommand.DisableAutoGenTag = true
			if err := doc.GenMarkdownTree(facadeCommand, usageDocPath); err != nil {
				return fmt.Errorf("generating usage documents: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&manPagePath, "man", "", "./doc/man", "A path to save man pages")
	cmd.Flags().StringVarP(&usageDocPath, "usage", "", "./doc/usage", "A path to save markdown usages")
	return cmd, nil
}
