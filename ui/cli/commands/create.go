package commands

import (
	"context"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/create"
	"github.com/kyoh86/gogh/v3/app/create_from_template"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
)

func NewCreateCommand(svc *ServiceSet) *cobra.Command {
	createUseCase := create.NewUseCase(
		svc.hostingService,
		svc.workspaceService,
		svc.gitService,
	)
	createFromTemplateUseCase := create_from_template.NewUseCase(
		svc.hostingService,
		svc.workspaceService,
		svc.gitService,
	)

	// TODO: Split flags from defaults
	// - Load defaults from config to config.CreateFlags (except that the subcommand is "man")
	// - If the flag is not set in config, use the default value (e.g. config.Create.CloneRetryLimit == 0 => 5)
	// - If the flag is set in config, use the value from config (e.g. config.Create.CloneRetryLimit == 5 => 5)
	// - If the flag is set in command line, use the value from command line (e.g. --clone-retry-limit 10 => 10)
	// ref: infra/config/token_store.go depends on core/auth/token_service.go
	var f config.CreateFlags

	checkFlags := func(_ context.Context, args []string) (*repository.ReferenceWithAlias, error) {
		var name string
		if len(args) == 0 {
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().
					Title("A ref of repository name to create").
					Validate(func(s string) error {
						_, err := svc.referenceParser.Parse(s)
						return err
					}).
					Value(&name),
			)).Run(); err != nil {
				return nil, err
			}
			return svc.referenceParser.ParseWithAlias(name)
		}
		return svc.referenceParser.ParseWithAlias(args[0])
	}

	runFunc := func(ctx context.Context, ref *repository.ReferenceWithAlias) error {
		if f.Template == "" {
			ropt := create.Options{
				TryCloneNotify: service.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				CreateRepositoryOptions: hosting.CreateRepositoryOptions{
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
				Alias: ref.Alias,
			}
			if err := createUseCase.Execute(ctx, ref.Reference, ropt); err != nil {
				return err
			}
		} else {
			template, err := svc.referenceParser.Parse(f.Template)
			if err != nil {
				return err
			}
			if err := createFromTemplateUseCase.Execute(ctx, ref.Reference, *template, create_from_template.CreateFromTemplateOptions{
				TryCloneNotify: service.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
				CreateRepositoryFromTemplateOptions: hosting.CreateRepositoryFromTemplateOptions{
					Description:        f.Description,
					IncludeAllBranches: f.IncludeAllBranches,
					Private:            f.Private,
				},
				Alias: ref.Alias,
			}); err != nil {
				return err
			}
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "create [flags] [[OWNER/]NAME]",
		Aliases: []string{"new"},
		Short:   "Create a new local and remote repository",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ref, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}
			if err := runFunc(ctx, ref); err != nil {
				log.FromContext(ctx).Errorf("failed to create repository: %v", err)
			}
			return nil
		},
	}
	// TODO: Validate flag combinations
	cmd.Flags().
		BoolVarP(&f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().
		StringVarP(&f.Template, "template", "", svc.defaults.Create.Template, "Create new repository from the template")
	cmd.Flags().
		BoolVarP(&f.IncludeAllBranches, "include-all-branches", "", svc.defaults.Create.IncludeAllBranches, "Create all branches in the template")
	cmd.Flags().
		StringVarP(&f.Description, "description", "", "", "A short description of the repository")
	cmd.Flags().
		StringVarP(&f.Homepage, "homepage", "", "", "A URL with more information about the repository")
	cmd.Flags().
		StringVarP(&f.LicenseTemplate, "license-template", "", svc.defaults.Create.LicenseTemplate, `Choose an open source license template that best suits your needs, and then use the license keyword as the license_template string when "auto-init" flag is set. For example, "mit" or "mpl-2.0"`)
	cmd.Flags().
		StringVarP(&f.GitignoreTemplate, "gitignore-template", "", svc.defaults.Create.GitignoreTemplate, `Desired language or platform .gitignore template to apply when "auto-init" flag is set. Use the name of the template without the extension. For example, "Haskell"`)
	cmd.Flags().
		BoolVarP(&f.Private, "private", "", svc.defaults.Create.Private, "Whether the repository is private")
	cmd.Flags().
		BoolVarP(&f.IsTemplate, "is-template", "", false, "Whether the repository is available as a template")
	cmd.Flags().
		BoolVarP(&f.DisableDownloads, "disable-downloads", "", svc.defaults.Create.DisableDownloads, `Disable "Downloads" page`)
	cmd.Flags().
		BoolVarP(&f.DisableWiki, "disable-wiki", "", svc.defaults.Create.DisableWiki, `Disable Wiki for the repository`)
	cmd.Flags().
		BoolVarP(&f.AutoInit, "auto-init", "", svc.defaults.Create.AutoInit, "Create an initial commit with empty README")
	cmd.Flags().
		BoolVarP(&f.DisableProjects, "disable-projects", "", svc.defaults.Create.DisableProjects, `Disable projects for the repository`)
	cmd.Flags().
		BoolVarP(&f.DisableIssues, "disable-issues", "", svc.defaults.Create.DisableIssues, `Disable issues for the repository`)
	cmd.Flags().
		BoolVarP(&f.PreventSquashMerge, "prevent-squash-merge", "", svc.defaults.Create.PreventSquashMerge, "Prevent squash-merging pull requests")
	cmd.Flags().
		BoolVarP(&f.PreventMergeCommit, "prevent-merge-commit", "", svc.defaults.Create.PreventMergeCommit, "Prevent merging pull requests with a merge commit")
	cmd.Flags().
		BoolVarP(&f.PreventRebaseMerge, "prevent-rebase-merge", "", svc.defaults.Create.PreventRebaseMerge, "Prevent rebase-merging pull requests")
	cmd.Flags().
		BoolVarP(&f.DeleteBranchOnMerge, "delete-branch-on-merge", "", svc.defaults.Create.DeleteBranchOnMerge, "Allow automatically deleting head branches when pull requests are merged")
	cmd.Flags().
		IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.defaults.Create.CloneRetryLimit, "")
	return cmd
}
