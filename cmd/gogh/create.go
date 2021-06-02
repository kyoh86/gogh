package main

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

type createFlagsStruct struct {
	Template            string `yaml:"template,omitempty"`
	Description         string `yaml:"-"`
	Homepage            string `yaml:"-"`
	LicenseTemplate     string `yaml:"licenseTemplate,omitempty"`
	GitignoreTemplate   string `yaml:"gitignoreTemplate,omitempty"`
	Private             bool   `yaml:"private,omitempty"`
	IsTemplate          bool   `yaml:"-"`
	DisableDownloads    bool   `yaml:"disableDownloads,omitempty"`
	DisableWiki         bool   `yaml:"disableWiki,omitempty"`
	AutoInit            bool   `yaml:"autoInit,omitempty"`
	DisableProjects     bool   `yaml:"disableProjects,omitempty"`
	DisableIssues       bool   `yaml:"disableIssues,omitempty"`
	PreventSquashMerge  bool   `yaml:"preventSquashMerge,omitempty"`
	PreventMergeCommit  bool   `yaml:"preventMergeCommit,omitempty"`
	PreventRebaseMerge  bool   `yaml:"preventRebaseMerge,omitempty"`
	DeleteBranchOnMerge bool   `yaml:"deleteBranchOnMerge,omitempty"`
	Dryrun              bool   `yaml:"-"`
}

var createFlags createFlagsStruct

var createCommand = &cobra.Command{
	Use:     "create [flags] [[OWNER/]NAME]",
	Aliases: []string{"new"},
	Short:   "Create a new project with a remote repository",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		var name string
		parser := gogh.NewSpecParser(&servers)
		if len(specs) == 0 {
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
		spec, server, err := parser.Parse(name)
		if err != nil {
			return err
		}

		local := gogh.NewLocalController(defaultRoot())
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

		if createFlags.Template == "" {
			ropt := &gogh.RemoteCreateOption{
				Description:         createFlags.Description,
				Homepage:            createFlags.Homepage,
				LicenseTemplate:     createFlags.LicenseTemplate,
				GitignoreTemplate:   createFlags.GitignoreTemplate,
				Private:             createFlags.Private,
				IsTemplate:          createFlags.IsTemplate,
				DisableDownloads:    createFlags.DisableDownloads,
				DisableWiki:         createFlags.DisableWiki,
				AutoInit:            createFlags.AutoInit,
				DisableProjects:     createFlags.DisableProjects,
				DisableIssues:       createFlags.DisableIssues,
				PreventSquashMerge:  createFlags.PreventSquashMerge,
				PreventMergeCommit:  createFlags.PreventMergeCommit,
				PreventRebaseMerge:  createFlags.PreventRebaseMerge,
				DeleteBranchOnMerge: createFlags.DeleteBranchOnMerge,
			}
			if server.User() != spec.Owner() {
				ropt.Organization = spec.Owner()
			}

			_, err = remote.Create(ctx, spec.Name(), ropt)
			return err
		}

		from, err := gogh.ParseSiblingSpec(spec, createFlags.Template)
		if err != nil {
			return err
		}
		ropt := &gogh.RemoteCreateFromTemplateOption{}
		if server.User() != spec.Owner() {
			ropt.Owner = spec.Owner()
		}
		if createFlags.Private {
			ropt.Private = true
		}
		_, err = remote.CreateFromTemplate(ctx, from.Owner(), from.Name(), spec.Name(), ropt)
		return err
	},
}

func init() {
	setup()
	createCommand.Flags().BoolVarP(&createFlags.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	createCommand.Flags().StringVarP(&createFlags.Template, "template", "", config.DefaultFlag.Create.Template, "Create new repository from the template")
	createCommand.Flags().StringVarP(&createFlags.Description, "description", "", "", "A short description of the repository")
	createCommand.Flags().StringVarP(&createFlags.Homepage, "homepage", "", "", "A URL with more information about the repository")
	createCommand.Flags().StringVarP(&createFlags.LicenseTemplate, "license-template", "", config.DefaultFlag.Create.LicenseTemplate, `Choose an open source license template that best suits your needs, and then use the license keyword as the license_template string when "auto-init" flag is set. For example, "mit" or "mpl-2.0"`)
	createCommand.Flags().StringVarP(&createFlags.GitignoreTemplate, "gitignore-template", "", config.DefaultFlag.Create.GitignoreTemplate, `Desired language or platform .gitignore template to apply when "auto-init" flag is set. Use the name of the template without the extension. For example, "Haskell"`)
	createCommand.Flags().BoolVarP(&createFlags.Private, "private", "", config.DefaultFlag.Create.Private, "Whether the repository is private")
	createCommand.Flags().BoolVarP(&createFlags.IsTemplate, "is-template", "", false, "Whether the repository is available as a template")
	createCommand.Flags().BoolVarP(&createFlags.DisableDownloads, "disable-downloads", "", config.DefaultFlag.Create.DisableDownloads, `Disable "Downloads" page`)
	createCommand.Flags().BoolVarP(&createFlags.DisableWiki, "disable-wiki", "", config.DefaultFlag.Create.DisableWiki, `Disable Wiki for the repository`)
	createCommand.Flags().BoolVarP(&createFlags.AutoInit, "auto-init", "", config.DefaultFlag.Create.AutoInit, "Create an initial commit with empty README")
	createCommand.Flags().BoolVarP(&createFlags.DisableProjects, "disable-projects", "", config.DefaultFlag.Create.DisableProjects, `Disable projects for the repository`)
	createCommand.Flags().BoolVarP(&createFlags.DisableIssues, "disable-issues", "", config.DefaultFlag.Create.DisableIssues, `Disable issues for the repository`)
	createCommand.Flags().BoolVarP(&createFlags.PreventSquashMerge, "prevent-squash-merge", "", config.DefaultFlag.Create.PreventSquashMerge, "Prevent squash-merging pull requests")
	createCommand.Flags().BoolVarP(&createFlags.PreventMergeCommit, "prevent-merge-commit", "", config.DefaultFlag.Create.PreventMergeCommit, "Prevent merging pull requests with a merge commit")
	createCommand.Flags().BoolVarP(&createFlags.PreventRebaseMerge, "prevent-rebase-merge", "", config.DefaultFlag.Create.PreventRebaseMerge, "Prevent rebase-merging pull requests")
	createCommand.Flags().BoolVarP(&createFlags.DeleteBranchOnMerge, "delete-branch-on-merge", "", config.DefaultFlag.Create.DeleteBranchOnMerge, "Allow automatically deleting head branches when pull requests are merged")
	facadeCommand.AddCommand(createCommand)
}
