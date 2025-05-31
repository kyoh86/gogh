package commands

import (
	"context"
	"fmt"

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
		Short:   "List overlay patterns",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			useCase := overlay_list.NewUseCase(svc.OverlayService)
			patterns, err := useCase.Execute(ctx)
			if err != nil {
				logger.Errorf("Failed to list overlay patterns: %v", err)
				return nil
			}

			for _, pattern := range patterns {
				fmt.Printf("Pattern: %s\n", pattern.Pattern)
				for _, file := range pattern.Files {
					fmt.Printf("  %s -> %s\n", file.SourcePath, file.TargetPath)
				}
				fmt.Println()
			}

			return nil
		},
	}
	return cmd, nil
}

func NewOverlayAddCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "add <pattern> <source-path> <target-path>",
		Short: "Add an overlay file pattern",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			pattern := args[0]
			sourcePath := args[1]
			targetPath := args[2]

			useCase := overlay_add.NewUseCase(svc.OverlayService)
			if err := useCase.Execute(ctx, pattern, sourcePath, targetPath); err != nil {
				return err
			}

			logger.Infof("Added overlay file %s -> %s for pattern %s", sourcePath, targetPath, pattern)
			return nil
		},
	}
	return cmd, nil
}

func NewOverlayRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "remove <pattern>",
		Aliases: []string{"rm"},
		Short:   "Remove an overlay pattern",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			pattern := args[0]

			useCase := overlay_remove.NewUseCase(svc.OverlayService)
			if err := useCase.Execute(ctx, pattern); err != nil {
				return err
			}

			logger.Infof("Removed overlay pattern %s", pattern)

			return nil
		},
	}
	return cmd, nil
}
