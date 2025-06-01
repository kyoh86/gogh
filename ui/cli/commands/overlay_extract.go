package commands

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/app/overlay_extract"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewOverlayExtractCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		pattern string
		force   bool
	}

	reposUseCase := repos.NewUseCase(svc.HostingService)

	checkFlags := func(ctx context.Context, args []string) (string, error) {
		if len(args) != 0 {
			return args[0], nil
		}
		var opts []huh.Option[string]
		for repo, err := range reposUseCase.Execute(ctx, repos.Options{}) {
			if err != nil {
				return "", fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		var selected string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("A repository to extract overlays").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return "", err
		}
		return selected, nil
	}

	cmd := &cobra.Command{
		Use:   "extract [repo-ref]",
		Short: "Extract untracked files as overlays",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := context.Background()
			selected, err := checkFlags(ctx, refs)
			if err != nil {
				return err
			}
			patternToUse := f.pattern

			overlayExtractUseCase := overlay_extract.NewUseCase(
				svc.GitService,
				svc.OverlayService,
				svc.WorkspaceService,
				svc.FinderService,
				svc.ReferenceParser,
			)

			overlayAddUseCase := overlay_add.NewUseCase(
				svc.OverlayService,
			)

			// Extract untracked files
			for result, err := range overlayExtractUseCase.Execute(ctx, selected, overlay_extract.Options{
				Pattern: f.pattern,
			}) {
				if err != nil {
					return err
				}

				// Determine pattern to use
				if patternToUse == "" {
					patternToUse = result.Reference.String()
				}

				if !f.force {
					var confirm bool
					if err := huh.NewForm(huh.NewGroup(
						huh.NewConfirm().
							Title(fmt.Sprintf("Are you sure you extract this file?\n%q", result.FilePath)).
							Value(&confirm),
					)).Run(); err != nil {
						return err
					}
					if !confirm {
						continue
					}
				}

				if err := overlayAddUseCase.Execute(ctx, result.FilePath, patternToUse, result.Content); err != nil {
					return fmt.Errorf("failed to register overlay for %s: %w", result.FilePath, err)
				}
				fmt.Printf("Registered %s as overlay\n", result.FilePath)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&f.pattern, "pattern", "", "", "Custom pattern for overlay (default: repository reference)")
	cmd.Flags().BoolVarP(&f.force, "force", "", false, "Do NOT confirm to delete.")
	return cmd, nil
}
