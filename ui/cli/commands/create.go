package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/create"
	"github.com/kyoh86/gogh/v4/app/create_from_template"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewCreateCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	createUseCase := create.NewUseCase(
		svc.HostingService,
		svc.WorkspaceService,
		svc.OverlayService,
		svc.ReferenceParser,
		svc.GitService,
	)
	createFromTemplateUseCase := create_from_template.NewUseCase(
		svc.HostingService,
		svc.WorkspaceService,
		svc.OverlayService,
		svc.ReferenceParser,
		svc.GitService,
	)

	var f config.CreateFlags

	checkFlags := func(_ context.Context, args []string) (string, error) {
		if len(args) > 0 {
			return args[0], nil
		}
		var name string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				Title("A ref of repository name to create [OWNER/]NAME[=ALIAS]").
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
				TryCloneOptions: try_clone.Options{
					Notify: try_clone.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
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
			if err := createUseCase.Execute(ctx, refWithAlias, ropt); err != nil {
				return fmt.Errorf("creating the repository: %w", err)
			}
		} else {
			template, err := svc.ReferenceParser.Parse(f.Template)
			if err != nil {
				return fmt.Errorf("invalid template: %w", err)
			}
			if err := createFromTemplateUseCase.Execute(ctx, refWithAlias, *template, create_from_template.CreateFromTemplateOptions{
				TryCloneOptions: try_clone.Options{
					Timeout: f.CloneRetryTimeout,
					Notify:  try_clone.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				},
				RepositoryOptions: create_from_template.RepositoryOptions{
					Description:        f.Description,
					IncludeAllBranches: f.IncludeAllBranches,
					Private:            f.Private,
				},
			}); err != nil {
				return fmt.Errorf("creating the repository from template: %w", err)
			}
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "create [flags] [[OWNER/]NAME[=ALIAS]]",
		Aliases: []string{"new"},
		Short:   "Create a new local and remote repository",
		Args:    cobra.RangeArgs(0, 1),
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
	flags.BoolVarP(cmd, &f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().Lookup("dryrun").NoOptDefVal = "false"
	cmd.Flags().StringVarP(&f.Template, "template", "", svc.Flags.Create.Template, "Create new repository from the template")
	flags.BoolVarP(cmd, &f.IncludeAllBranches, "include-all-branches", "", svc.Flags.Create.IncludeAllBranches, "Create all branches in the template")
	cmd.Flags().Lookup("include-all-branches").NoOptDefVal = "false"
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
	flags.BoolVarP(cmd, &f.Private, "private", "", svc.Flags.Create.Private, "Whether the repository is private")
	cmd.Flags().Lookup("private").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.IsTemplate, "is-template", "", false, "Whether the repository is available as a template")
	cmd.Flags().Lookup("is-template").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.DisableDownloads, "disable-downloads", "", svc.Flags.Create.DisableDownloads, `Disable "Downloads" page`)
	cmd.Flags().Lookup("disable-downloads").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.DisableWiki, "disable-wiki", "", svc.Flags.Create.DisableWiki, `Disable Wiki for the repository`)
	cmd.Flags().Lookup("disable-wiki").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.AutoInit, "auto-init", "", svc.Flags.Create.AutoInit, "Create an initial commit with empty README")
	cmd.Flags().Lookup("disable-wiki").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.DisableProjects, "disable-projects", "", svc.Flags.Create.DisableProjects, `Disable projects for the repository`)
	flags.BoolVarP(cmd, &f.DisableIssues, "disable-issues", "", svc.Flags.Create.DisableIssues, `Disable issues for the repository`)
	flags.BoolVarP(cmd, &f.PreventSquashMerge, "prevent-squash-merge", "", svc.Flags.Create.PreventSquashMerge, "Prevent squash-merging pull requests")
	flags.BoolVarP(cmd, &f.PreventMergeCommit, "prevent-merge-commit", "", svc.Flags.Create.PreventMergeCommit, "Prevent merging pull requests with a merge commit")
	flags.BoolVarP(cmd, &f.PreventRebaseMerge, "prevent-rebase-merge", "", svc.Flags.Create.PreventRebaseMerge, "Prevent rebase-merging pull requests")
	flags.BoolVarP(cmd, &f.DeleteBranchOnMerge, "delete-branch-on-merge", "", svc.Flags.Create.DeleteBranchOnMerge, "Allow automatically deleting head branches when pull requests are merged")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Create.CloneRetryTimeout, "Timeout for each clone attempt")
	cmd.Flags().IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Create.CloneRetryLimit, "The number of retries to clone a repository")
	return cmd, nil
}
