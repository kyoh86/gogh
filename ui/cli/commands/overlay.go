package commands

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/app/overlay_list"
	"github.com/kyoh86/gogh/v4/app/overlay_remove"
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

func NewOverlayListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		Short:   "List overlays",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entries, err := overlay_list.NewUseCase(svc.OverlayStore).Execute(ctx)
			if err != nil {
				return fmt.Errorf("listing overlay: %w", err)
			}

			for _, entry := range entries {
				fmt.Printf("- ")
				fmt.Printf("Repository pattern: %s\n", entry.RepoPattern)
				if entry.ForInit {
					fmt.Printf("  For Init: Yes\n")
				}
				fmt.Printf("  Overlay path: %s\n", entry.RelativePath)
			}

			return nil
		},
	}
	return cmd, nil
}

func NewOverlayAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		forInit bool
	}
	cmd := &cobra.Command{
		Use:   "add [flags] <source-path> <repo-pattern> <target-path>",
		Short: "Add an overlay file",
		Args:  cobra.ExactArgs(3),
		Example: `   Add an overlay file to a repository.
   The <source-path> is the path to the file you want to add as an overlay.
   The <repo-pattern> is the pattern of the repository you want to add the overlay to.
   The <target-path> is the path where the overlay file will be copied to in the repository.

   For example, to add a custom VSCode settings file to a repository, you can run:

     gogh overlay add /path/to/source/vscode/settings.json "github.com/owner/repo" .vscode/settings.json

   The overlay file will be copied to the repository when you run ` + "`gogh create`, `gogh clone` or `gogh fork`" + `.

   You can also apply template files only for the ` + "`gogh create`" + ` command by using the ` + "`--for-init`" + ` flag:

     gogh overlay add --for-init /path/to/source/deno.jsonc "github.com/owner/deno-*" deno.jsonc

   This will copy the ` + "`deno.jsonc`" + ` file to the root of the repository only when you run ` + "`gogh create`" + `
   if the repository matches the pattern ` + "`github.com/owner/deno-*`" + `.

   You can create overlay files that never be applied to the repository automatically,
   (and only be applied manually by ` + "`gogh overlay apply`" + ` command),
   you can set the ` + "`--repo-pattern`" + ` flag to never match any repository.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			sourcePath := args[0]
			repoPattern := args[1]
			targetPath := args[2]
			if filepath.IsAbs(targetPath) {
				return fmt.Errorf("target path must be relative, got absolute path: %s", targetPath)
			}

			if err := overlay_add.NewUseCase(svc.OverlayStore).Execute(ctx, f.forInit, repoPattern, targetPath, sourcePath); err != nil {
				return err
			}

			logger.Infof("Added overlay file %s -> %s for repo-pattern %s", sourcePath, targetPath, repoPattern)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.forInit, "for-init", "", false, "Register the overlay for 'gogh create' command")
	return cmd, nil
}

func NewOverlayRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		forInit bool
	}
	cmd := &cobra.Command{
		Use:     "remove [flags] <repo-pattern> <target-path>",
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove an overlay",
		Args:    cobra.ExactArgs(2),
		Example: `   Remove an overlay file from a repository.
			 The <repo-pattern> is the pattern of the repository you want to remove the overlay from.
			 The <target-path> is the path where the overlay file is located in the repository.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			repoPattern := args[0]
			targetPath := args[1]
			if filepath.IsAbs(targetPath) {
				return fmt.Errorf("target path must be relative, got absolute path: %s", targetPath)
			}

			if err := overlay_remove.NewUseCase(svc.OverlayStore).Execute(ctx, f.forInit, targetPath, repoPattern); err != nil {
				return err
			}

			logger.Infof("Removed overlay %s for %s", targetPath, repoPattern)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.forInit, "for-init", "", false, "Remove the overlay for 'gogh create' command")
	return cmd, nil
}
