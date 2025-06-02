package commands

import (
	"context"
	"fmt"
	"os"
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
	}
	return cmd, nil
}

func NewOverlayListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List overlays",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			useCase := overlay_list.NewUseCase(svc.OverlayService)
			entries, err := useCase.Execute(ctx)
			if err != nil {
				return fmt.Errorf("listing overlay: %w", err)
			}

			for _, entry := range entries {
				fmt.Printf("- ")
				fmt.Printf("Repository pattern: %s\n", entry.Pattern)
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
		pattern string
		forInit bool
	}
	cmd := &cobra.Command{
		Use:   "add [options] <source-path> <target-path>",
		Short: "Add an overlay file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			sourcePath := args[0]
			targetPath := args[1]
			if filepath.IsAbs(targetPath) {
				return fmt.Errorf("target path must be relative, got absolute path: %s", targetPath)
			}

			source, err := os.Open(sourcePath)
			if err != nil {
				return fmt.Errorf("opening source file %s: %w", sourcePath, err)
			}
			defer source.Close()

			useCase := overlay_add.NewUseCase(svc.OverlayService)
			if err := useCase.Execute(ctx, f.forInit, f.pattern, targetPath, source); err != nil {
				return err
			}

			logger.Infof("Added overlay file %s -> %s for pattern %s", sourcePath, targetPath, f.pattern)
			return nil
		},
	}
	cmd.Flags().StringVarP(&f.pattern, "pattern", "p", "", "Pattern to match repositories (e.g., 'github.com/owner/repo', '**/gogh')")
	cmd.Flags().BoolVar(&f.forInit, "for-init", false, "Apply this overlay for 'gogh create' command")
	return cmd, nil
}

func NewOverlayRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		pattern string
		forInit bool
	}
	cmd := &cobra.Command{
		Use:     "remove <target-path>",
		Aliases: []string{"rm", "del", "delete"},
		Short:   "Remove an overlay pattern",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			targetPath := args[0]
			if filepath.IsAbs(targetPath) {
				return fmt.Errorf("target path must be relative, got absolute path: %s", targetPath)
			}

			useCase := overlay_remove.NewUseCase(svc.OverlayService)
			if err := useCase.Execute(ctx, f.forInit, targetPath, f.pattern); err != nil {
				return err
			}

			logger.Infof("Removed overlay %s for %s", targetPath, f.pattern)
			return nil
		},
	}
	cmd.Flags().StringVarP(&f.pattern, "pattern", "p", "", "Pattern to match repositories (e.g., 'github.com/owner/repo', '**/gogh')")
	cmd.Flags().BoolVar(&f.forInit, "for-init", false, "Remove this overlay for 'gogh create' command")
	return cmd, nil
}
