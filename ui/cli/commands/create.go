package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/clone/try"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/create"
	"github.com/kyoh86/gogh/v4/app/create/template"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewCreateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.CreateFlags

	checkFlags := func(_ context.Context, args []string) (string, error) {
		if len(args) > 0 {
			return args[0], nil
		}
		var name string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				Title("A ref of repository name to create [[<host>/]<owner>/]<name>[=<alias>]").
				Validate(func(s string) error {
					// Never do
					_, err := svc.ReferenceParser.Parse(s)
					return err
				}).
				Value(&name),
		)).Run(); err != nil {
			return "", err
		}
		return name, nil
	}

	runFunc := func(ctx context.Context, refWithAlias string) error {
		if f.Template == "" {
			ropt := create.Options{
				TryCloneOptions: try.Options{
					Notify: try.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				},
				RepositoryOptions: create.RepositoryOptions{
					Description:         f.Description,
					Homepage:            f.Homepage,
					LicenseTemplate:     f.LicenseTemplate,
					GitignoreTemplate:   f.GitignoreTemplate,
					Private:             f.Private,
					IsTemplate:          f.IsTemplate,
					DisableDownloads:    f.DisableDownloads,
					DisableWiki:         f.DisableWiki,
					AutoInit:            f.AutoInit,
					DisableProjects:     f.DisableProjects,
					DisableIssues:       f.DisableIssues,
					PreventSquashMerge:  f.PreventSquashMerge,
					PreventMergeCommit:  f.PreventMergeCommit,
					PreventRebaseMerge:  f.PreventRebaseMerge,
					DeleteBranchOnMerge: f.DeleteBranchOnMerge,
				},
			}
			if err := create.NewUsecase(
				svc.HostingService,
				svc.WorkspaceService,
				svc.FinderService,
				svc.OverlayService,
				svc.ScriptService,
				svc.HookService,
				svc.ReferenceParser,
				svc.GitService,
			).Execute(ctx, refWithAlias, ropt); err != nil {
				return fmt.Errorf("creating the repository: %w", err)
			}
		} else {
			tmp, err := svc.ReferenceParser.Parse(f.Template)
			if err != nil {
				return fmt.Errorf("invalid template: %w", err)
			}
			if err := template.NewUsecase(
				svc.HostingService,
				svc.WorkspaceService,
				svc.FinderService,
				svc.OverlayService,
				svc.ScriptService,
				svc.HookService,
				svc.ReferenceParser,
				svc.GitService,
			).Execute(ctx, refWithAlias, *tmp, template.CreateFromTemplateOptions{
				TryCloneOptions: try.Options{
					Timeout: f.CloneRetryTimeout,
					Notify:  try.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				},
				RepositoryOptions: template.RepositoryOptions{
					Description:        f.Description,
					IncludeAllBranches: f.IncludeAllBranches,
					Private:            f.Private,
				},
			}); err != nil {
				return fmt.Errorf("creating the repository from template: %w", err)
			}
		}
		if f.DryRun {
			fmt.Printf("Apply overlay for %q\n", refWithAlias)
			return nil
		}

		return nil
	}

	cmd := &cobra.Command{
		Use:     "create [flags] [[[<host>/]<owner>/]<name>[=<alias>]]",
		Aliases: []string{"new"},
		Short:   "Create a new local and remote repository",
		Args:    cobra.RangeArgs(0, 1),
		Example: `  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is used for a local repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ref, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}
			log.FromContext(ctx).Infof("Creating %q", ref)
			if err := runFunc(ctx, ref); err != nil {
				return err
			}
			log.FromContext(ctx).Infof("Created %q", ref)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.DryRun, "dry-run", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().StringVarP(&f.Template, "template", "", svc.Flags.Create.Template, "Create new repository from the template")
	cmd.Flags().BoolVarP(&f.IncludeAllBranches, "include-all-branches", "", svc.Flags.Create.IncludeAllBranches, "Create all branches in the template")
	cmd.Flags().StringVarP(&f.Description, "description", "", "", "A short description of the repository")
	cmd.Flags().StringVarP(&f.Homepage, "homepage", "", "", "A URL with more information about the repository")
	cmd.Flags().StringVarP(
		&f.LicenseTemplate,
		"license-template",
		"",
		svc.Flags.Create.LicenseTemplate,
		strings.Join([]string{
			`Choose an open source license template that best suits your needs,`,
			`and then use the license keyword as the license_template string when "auto-init" flag is set.`,
			`For example, "mit" or "mpl-2.0"`,
		}, " "),
	)
	cmd.Flags().StringVarP(&f.GitignoreTemplate, "gitignore-template", "", svc.Flags.Create.GitignoreTemplate, `Desired language or platform .gitignore template to apply when "auto-init" flag is set. Use the name of the template without the extension. For example, "Haskell"`)
	cmd.Flags().BoolVarP(&f.Private, "private", "", svc.Flags.Create.Private, "Whether the repository is private")
	cmd.Flags().BoolVarP(&f.IsTemplate, "is-template", "", false, "Whether the repository is available as a template")
	cmd.Flags().BoolVarP(&f.DisableDownloads, "disable-downloads", "", svc.Flags.Create.DisableDownloads, `Disable "Downloads" page`)
	cmd.Flags().BoolVarP(&f.DisableWiki, "disable-wiki", "", svc.Flags.Create.DisableWiki, `Disable Wiki for the repository`)
	cmd.Flags().BoolVarP(&f.AutoInit, "auto-init", "", svc.Flags.Create.AutoInit, "Create an initial commit with empty README")
	cmd.Flags().BoolVarP(&f.DisableProjects, "disable-projects", "", svc.Flags.Create.DisableProjects, `Disable projects for the repository`)
	cmd.Flags().BoolVarP(&f.DisableIssues, "disable-issues", "", svc.Flags.Create.DisableIssues, `Disable issues for the repository`)
	cmd.Flags().BoolVarP(&f.PreventSquashMerge, "prevent-squash-merge", "", svc.Flags.Create.PreventSquashMerge, "Prevent squash-merging pull requests")
	cmd.Flags().BoolVarP(&f.PreventMergeCommit, "prevent-merge-commit", "", svc.Flags.Create.PreventMergeCommit, "Prevent merging pull requests with a merge commit")
	cmd.Flags().BoolVarP(&f.PreventRebaseMerge, "prevent-rebase-merge", "", svc.Flags.Create.PreventRebaseMerge, "Prevent rebase-merging pull requests")
	cmd.Flags().BoolVarP(&f.DeleteBranchOnMerge, "delete-branch-on-merge", "", svc.Flags.Create.DeleteBranchOnMerge, "Allow automatically deleting head branches when pull requests are merged")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Create.CloneRetryTimeout, "Timeout for each clone attempt")
	cmd.Flags().IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Create.CloneRetryLimit, "The number of retries to clone a repository")
	return cmd, nil
}
