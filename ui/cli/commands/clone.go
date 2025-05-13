package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/clone"
	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCloneCommand(_ context.Context, svc *ServiceSet) *cobra.Command {
	var f struct {
		dryrun bool
	}

	reposUseCase := repos.NewUseCase(svc.hostingService)
	cloneUseCase := clone.NewUseCase(svc.hostingService, svc.workspaceService, svc.gitService)

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		var opts []huh.Option[string]
		for repo, err := range reposUseCase.Execute(ctx, repos.Options{}) {
			if err != nil {
				return nil, err
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("A remote repository to clone").
				Options(opts...).
				Value(&args),
		)).Run(); err != nil {
			return nil, err
		}
		return args, nil
	}

	preprocessFlags := func(ctx context.Context, args []string) ([]repository.ReferenceWithAlias, error) {
		args, err := checkFlags(ctx, args)
		if err != nil {
			return nil, err
		}
		refs := make([]repository.ReferenceWithAlias, 0, len(args))
		for _, s := range args {
			ref, err := svc.referenceParser.ParseWithAlias(s)
			if err != nil {
				return nil, err
			}
			refs = append(refs, *ref)
		}
		return refs, nil
	}

	runFunc := func(ctx context.Context, refs []repository.ReferenceWithAlias) error {
		if f.dryrun {
			for _, ref := range refs {
				if ref.Alias == nil {
					log.FromContext(ctx).Infof("git clone %q", ref.Reference)
				} else {
					log.FromContext(ctx).Infof("git clone %q into %q", ref.Reference, ref.Alias)
				}
			}
			return nil
		}

		eg, ctx := errgroup.WithContext(ctx)
		for _, ref := range refs {
			eg.Go(func() error {
				return cloneUseCase.Execute(ctx, ref.Reference, clone.Options{
					Alias:          ref.Alias,
					TryCloneNotify: service.RetryLimit(1, nil),
				})
			})
		}
		return eg.Wait()
	}

	c := &cobra.Command{
		Use:     "clone [flags] [[OWNER/]NAME[=ALIAS]]...",
		Aliases: []string{"get"},
		Short:   "Clone remote repositories to local",
		Example: `  It accepts a short notation for a remote repository
  (for example, "github.com/kyoh86/example") like below.
    - "NAME": e.g. "example"; 
    - "OWNER/NAME": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default".

  It accepts an alias for each repository.
	The alias is a local name for the remote repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example"
  into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			refs, err := preprocessFlags(ctx, args)
			if err != nil {
				return err
			}
			if err := runFunc(ctx, refs); err != nil {
				log.FromContext(ctx).Errorf("failed to clone repositories: %v", err)
			}
			return nil
		},
	}

	c.Flags().
		BoolVarP(&f.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	return c
}
