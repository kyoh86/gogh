package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/app/overlay_extract"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewOverlayExtractCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		repoPattern string
		forInit     bool
		force       bool
	}

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		var opts []huh.Option[string]
		for repo, err := range repos.NewUseCase(svc.HostingService).Execute(ctx, repos.Options{}) {
			if err != nil {
				return nil, fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		var selected []string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Repositories to extract overlays").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return nil, err
		}
		return selected, nil
	}

	cmd := &cobra.Command{
		Use:   "extract [flags] [[[<host>/]<owner>/]<name>...]",
		Short: "Extract untracked files as overlays",
		Args:  cobra.ArbitraryArgs,
		Example: `  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is a local name for the remote repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			refs, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}
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
			for _, ref := range refs {
				logger.Infof("Extracting files from %q", ref)
				if err := view.ProcessWithConfirmation(ctx, overlayExtractUseCase.Execute(ctx, ref, overlay_extract.Options{}),
					func(result *overlay_extract.ExtractResult) string {
						return fmt.Sprintf("Extract %q from %q", result.FilePath, ref)
					},
					func(result *overlay_extract.ExtractResult) error {
						repoPattern := f.repoPattern
						// Determine repo-pattern to use
						if repoPattern == "" {
							repoPattern = result.Reference.String()
						}
						if err := overlayAddUseCase.Execute(ctx, f.forInit, result.RelativePath, repoPattern, result.FilePath); err != nil {
							return fmt.Errorf("failed to register overlay for %s: %w", result.FilePath, err)
						}
						logger.Infof("Registered %q from %q as overlay\n", result.FilePath, ref)
						return nil
					},
				); err != nil {
					if errors.Is(err, view.ErrQuit) {
						return nil
					}
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&f.repoPattern, "repo-pattern", "p", "", "Pattern to match repositories to apply the overlays to (e.g., 'github.com/owner/repo', '**/gogh'; default: repository reference)")
	cmd.Flags().BoolVarP(&f.force, "force", "", false, "Do NOT confirm to extract for each file")
	cmd.Flags().BoolVarP(&f.forInit, "for-init", "", false, "Register the overlay for `gogh create` command")
	return cmd, nil
}
