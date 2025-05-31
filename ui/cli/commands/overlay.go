package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/workspace"
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

			patterns := svc.OverlayService.GetPatterns()
			if len(patterns) == 0 {
				logger.Info("No overlay patterns defined")
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

			patterns := svc.OverlayService.GetPatterns()
			var files []workspace.OverlayFile

			// Find existing pattern
			for _, p := range patterns {
				if p.Pattern == pattern {
					files = p.Files
					break
				}
			}

			// Add new file
			files = append(files, workspace.OverlayFile{
				SourcePath: sourcePath,
				TargetPath: targetPath,
			})

			if err := svc.OverlayService.AddPattern(pattern, files); err != nil {
				return fmt.Errorf("adding pattern %s: %w", pattern, err)
			}

			logger.Infof("Added overlay file %s -> %s for pattern %s", sourcePath, targetPath, pattern)
			svc.OverlayService.MarkSaved()

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

			if err := svc.OverlayService.RemovePattern(pattern); err != nil {
				return fmt.Errorf("removing pattern %s: %w", pattern, err)
			}

			logger.Infof("Removed overlay pattern %s", pattern)
			svc.OverlayService.MarkSaved()

			return nil
		},
	}
	return cmd, nil
}
