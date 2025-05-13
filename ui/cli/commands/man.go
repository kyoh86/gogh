package commands

import (
	"os"

	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewManCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "man",
		Short:  "Generate manual",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			facadeCommand := cmd.Parent()
			list, _, err := facadeCommand.Traverse([]string{"list"})
			if err != nil {
				return err
			}
			cwd, _, err := facadeCommand.Traverse([]string{"cwd"})
			if err != nil {
				return err
			}
			list.Flag("format").Usage = flags.LocationFormatLongUsage
			cwd.Flag("format").Usage = flags.LocationFormatLongUsage
			header := &doc.GenManHeader{
				Title:   "GOGH",
				Section: "1",
			}
			if err := os.MkdirAll("./doc/man", 0755); err != nil {
				return err
			}
			if err := doc.GenManTree(facadeCommand, header, "./doc/man"); err != nil {
				return err
			}
			if err := os.MkdirAll("./doc/usage", 0755); err != nil {
				return err
			}
			facadeCommand.DisableAutoGenTag = true
			if err := doc.GenMarkdownTree(facadeCommand, "./doc/usage"); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
