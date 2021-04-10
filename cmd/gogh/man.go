// +build man

package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const formatLongUsage = `
Print each project in a given format, where [format] can be one of "rel-path", "rel-file-path",
"full-file-path", "url", "fields" and "fields:[separator]".

- rel-path:

	A part of the URL to specify a repository.  For example: "github.com/kyoh86/gogh"

- rel-file-path:

	A relative file path of the project from gogh roots.  For example in windows:
	"github.com\kyoh86\gogh"; in other case: "github.com/kyoh86/gogh".

- full-file-path

	A full file path of the project.  For example in Windows:
	"C:\Users\kyoh86\Projects\github.com\kyoh86\gogh"; in other case:
	"/root/Projects/github.com/kyoh86/gogh".

- url

	A URL of the repository.

- fields

	Tab separated all formats and properties of the project.
	i.e. [full-file-path]\t[rel-file-path]\t[url]\t[rel-path]\t[host]\t[owner]\t[name]

- fields:[separator]

	Like "fields" but with the explicit separator.
`

var manFlags struct{}

var manCommand = &cobra.Command{
	Use:    "man",
	Short:  "Generate manual",
	Hidden: true,
	Args:   cobra.ExactArgs(0),
	RunE: func(*cobra.Command, []string) error {
		listCommand.Flag("format").Usage = formatLongUsage
		header := &doc.GenManHeader{
			Title:   "GOGH",
			Section: "1",
		}
		if err := doc.GenManTree(facadeCommand, header, "."); err != nil {
			return err
		}
		if err := os.MkdirAll("./usage", 0755); err != nil {
			return err
		}
		if err := doc.GenMarkdownTree(facadeCommand, "./usage"); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	facadeCommand.AddCommand(manCommand)
}
