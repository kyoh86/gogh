package main

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var createFlags struct {
	template            string
	description         string
	homepage            string
	licenseTemplate     string
	gitignoreTemplate   string
	private             bool
	isTemplate          bool
	disableDownloads    bool
	disableWiki         bool
	autoInit            bool
	disableProjects     bool
	disableIssues       bool
	preventSquashMerge  bool
	preventMergeCommit  bool
	preventRebaseMerge  bool
	deleteBranchOnMerge bool
	dryrun              bool
}

var createCommand = &cobra.Command{
	Use:     "create [flags] [[OWNER/]NAME]",
	Aliases: []string{"new"},
	Short:   "Create a new project with a remote repository",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		servers := Servers()
		var name string
		if len(specs) == 0 {
			parser := gogh.NewSpecParser(servers)
			if err := survey.AskOne(&survey.Input{
				Message: "A spec of repository name to create",
			}, &name, survey.WithValidator(func(input interface{}) error {
				s, ok := input.(string)
				if !ok {
					return errors.New("invalid type")
				}
				_, _, err := parser.Parse(s)
				return err
			})); err != nil {
				return err
			}
		} else {
			name = specs[0]
		}

		ctx := cmd.Context()
		parser := gogh.NewSpecParser(servers)
		spec, server, err := parser.Parse(name)
		if err != nil {
			return err
		}

		local := gogh.NewLocalController(DefaultRoot())
		if _, err = local.Create(ctx, spec, nil); err != nil {
			if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
				return err
			}
		}

		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		remote := gogh.NewRemoteController(adaptor)

		// check repo has already existed
		if _, err := remote.Get(ctx, spec.Owner(), spec.Name(), nil); err == nil {
			return nil
		}

		if createFlags.template == "" {
			ropt := &gogh.RemoteCreateOption{
				Description:         createFlags.description,
				Homepage:            createFlags.homepage,
				LicenseTemplate:     createFlags.licenseTemplate,
				GitignoreTemplate:   createFlags.gitignoreTemplate,
				Private:             createFlags.private,
				IsTemplate:          createFlags.isTemplate,
				DisableDownloads:    createFlags.disableDownloads,
				DisableWiki:         createFlags.disableWiki,
				AutoInit:            createFlags.autoInit,
				DisableProjects:     createFlags.disableProjects,
				DisableIssues:       createFlags.disableIssues,
				PreventSquashMerge:  createFlags.preventSquashMerge,
				PreventMergeCommit:  createFlags.preventMergeCommit,
				PreventRebaseMerge:  createFlags.preventRebaseMerge,
				DeleteBranchOnMerge: createFlags.deleteBranchOnMerge,
			}
			if server.User() != spec.Owner() {
				ropt.Organization = spec.Owner()
			}

			_, err = remote.Create(ctx, spec.Name(), ropt)
			return err
		}

		from, err := gogh.ParseSiblingSpec(spec, createFlags.template)
		if err != nil {
			return err
		}
		ropt := &gogh.RemoteCreateFromTemplateOption{}
		if server.User() != spec.Owner() {
			ropt.Owner = spec.Owner()
		}
		if createFlags.private {
			ropt.Private = true
		}
		_, err = remote.CreateFromTemplate(ctx, from.Owner(), from.Name(), spec.Name(), ropt)
		return err
	},
}

func init() {
	createCommand.Flags().BoolVarP(&createFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	createCommand.Flags().StringVarP(&createFlags.template, "template", "", Flag().Create.Template, "Create new repository from the template")
	createCommand.Flags().StringVarP(&createFlags.description, "description", "", "", "A short description of the repository")
	createCommand.Flags().StringVarP(&createFlags.homepage, "homepage", "", "", "A URL with more information about the repository")
	createCommand.Flags().StringVarP(&createFlags.licenseTemplate, "license-template", "", Flag().Create.LicenseTemplate, `Choose an open source license template that best suits your needs, and then use the license keyword as the license_template string when "auto-init" flag is set. For example, "mit" or "mpl-2.0"`)
	createCommand.Flags().StringVarP(&createFlags.gitignoreTemplate, "gitignore-template", "", Flag().Create.GitignoreTemplate, `Desired language or platform .gitignore template to apply when "auto-init" flag is set. Use the name of the template without the extension. For example, "Haskell"`)
	createCommand.Flags().BoolVarP(&createFlags.private, "private", "", Flag().Create.Private, "Whether the repository is private")
	createCommand.Flags().BoolVarP(&createFlags.isTemplate, "is-template", "", false, "Whether the repository is available as a template")
	createCommand.Flags().BoolVarP(&createFlags.disableDownloads, "disable-downloads", "", Flag().Create.DisableDownloads, `Disable "Downloads" page`)
	createCommand.Flags().BoolVarP(&createFlags.disableWiki, "disable-wiki", "", Flag().Create.DisableWiki, `Disable Wiki for the repository`)
	createCommand.Flags().BoolVarP(&createFlags.autoInit, "auto-init", "", Flag().Create.AutoInit, "Create an initial commit with empty README")
	createCommand.Flags().BoolVarP(&createFlags.disableProjects, "disable-projects", "", Flag().Create.DisableProjects, `Disable projects for the repository`)
	createCommand.Flags().BoolVarP(&createFlags.disableIssues, "disable-issues", "", Flag().Create.DisableIssues, `Disable issues for the repository`)
	createCommand.Flags().BoolVarP(&createFlags.preventSquashMerge, "prevent-squash-merge", "", Flag().Create.PreventSquashMerge, "Prevent squash-merging pull requests")
	createCommand.Flags().BoolVarP(&createFlags.preventMergeCommit, "prevent-merge-commit", "", Flag().Create.PreventMergeCommit, "Prevent merging pull requests with a merge commit")
	createCommand.Flags().BoolVarP(&createFlags.preventRebaseMerge, "prevent-rebase-merge", "", Flag().Create.PreventRebaseMerge, "Prevent rebase-merging pull requests")
	createCommand.Flags().BoolVarP(&createFlags.deleteBranchOnMerge, "delete-branch-on-merge", "", Flag().Create.DeleteBranchOnMerge, "Allow automatically deleting head branches when pull requests are merged")
	facadeCommand.AddCommand(createCommand)
}
