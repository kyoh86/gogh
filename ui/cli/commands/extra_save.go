package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/extra/save"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewExtraSaveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		confirmMode string
	}

	cmd := &cobra.Command{
		Use:   "save <repository>",
		Short: "Save excluded files as auto-apply extra",
		Long: `Save files that are excluded by .gitignore as auto-apply extra.
These extra will be automatically applied when the repository is cloned.`,
		Args: cobra.ExactArgs(1),
		Example: `  save github.com/kyoh86/example
  save .  # Save from current directory repository

  It accepts a short notation for the repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)

			repoStr := args[0]

			// Use current directory if reference is "."
			if repoStr == "." {
				repo, err := cwd.NewUsecase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
				if err != nil {
					return fmt.Errorf("finding repository from current directory: %w", err)
				}
				repoStr = repo.Ref().String()
			}

			uc := save.NewUsecase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.GitService,
				svc.OverlayService,
				svc.ScriptService,
				svc.HookService,
				svc.ExtraService,
				svc.ReferenceParser,
			)

			// Get excluded files
			result, err := uc.GetExcludedFiles(ctx, repoStr)
			if err != nil {
				return err
			}

			if len(result.Files) == 0 {
				logger.Warn("No excluded files found")
				return nil
			}

			// Select files based on confirmation mode
			var selectedFiles []string
			switch f.confirmMode {
			case "none":
				logger.Infof("Skipping confirmation, saving all %d excluded files", len(result.Files))
				selectedFiles = result.Files
			case "iterative":
				selection, err := view.ConfirmFilesIterative(ctx, result.RepositoryPath, result.Files)
				if err != nil {
					if errors.Is(err, view.ErrQuit) {
						logger.Info("File selection cancelled")
						return nil
					}
					return fmt.Errorf("selecting files: %w", err)
				}
				if len(selection.Selected) == 0 {
					logger.Info("No files selected")
					return nil
				}
				selectedFiles = selection.Selected
			case "select", "":
				// Default to select mode
				selection, err := view.SelectFiles(ctx, result.RepositoryPath, result.Files)
				if err != nil {
					return fmt.Errorf("selecting files: %w", err)
				}
				if len(selection.Selected) == 0 {
					logger.Info("No files selected")
					return nil
				}
				selectedFiles = selection.Selected
			default:
				return fmt.Errorf("invalid confirm mode: %s (valid options: select, iterative, none)", f.confirmMode)
			}

			// Save selected files
			if err := uc.SaveFiles(ctx, repoStr, selectedFiles); err != nil {
				return err
			}

			logger.Infof("Saved auto-apply extra for %s", repoStr)
			return nil
		},
	}

	cmd.Flags().StringVar(&f.confirmMode, "confirm-mode", "select", "Confirmation mode: select (multi-select), iterative (one-by-one), none (skip confirmation)")

	return cmd, nil
}
