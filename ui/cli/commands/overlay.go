package commands

import (
	"context"

	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "overlay",
		Short: "Manage repository overlay files",
		Example: `   Overlay files are used to put custom files into repositories.
   They are useful to add files that are not tracked by the repository, such as editor configurations or scripts.

   For example, to add a custom VSCode settings file to a repository, you can run:

     gogh overlay add /path/to/source/vscode/settings.json "github.com/owner/repo" .vscode/settings.json

   Then when you run ` + "`gogh create`, `gogh clone` or `gogh fork`" + `, the files will be copied to the repository.

   You can also apply template files only for the ` + "`gogh create`" + ` command by using the ` + "`--for-init`" + ` flag:

     gogh overlay add --for-init /path/to/source/deno.jsonc "github.com/owner/deno-*" deno.jsonc

   This will copy the ` + "`deno.jsonc`" + ` file to the root of the repository only when you run ` + "`gogh create`" + `
   if the repository matches the pattern ` + "`github.com/owner/deno-*`" + `.

   And then you can use the ` + "`gogh overlay apply`" + ` command to apply the overlay files manually.

   You can create overlay files that never be applied to the repository automatically,
   (and only be applied manually by ` + "`gogh overlay apply`" + ` command),
   you can set the ` + "`--repo-pattern`" + ` flag to never match any repository.`,
	}
	return cmd, nil
}
