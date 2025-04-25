package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/kyoh86/gogh/v3/domain/local"
	"github.com/kyoh86/gogh/v3/domain/remote"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewCreateCommand(conf *config.ConfigStore, tokens *config.TokenStore, defaults *config.FlagStore) *cobra.Command {
	var f config.CreateFlags
	cmd := &cobra.Command{
		Use:     "create [flags] [[OWNER/]NAME]",
		Aliases: []string{"new"},
		Short:   "Create a new local and remote repository",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			var name string
			parser := reporef.NewRepoRefParser(tokens.GetDefaultKey())
			if len(refs) == 0 {
				if err := huh.NewForm(huh.NewGroup(
					huh.NewInput().
						Title("A ref of repository name to create").
						Validate(func(s string) error {
							_, err := parser.Parse(s)
							return err
						}).
						Value(&name),
				)).Run(); err != nil {
					return err
				}
			} else {
				name = refs[0]
			}

			ctx := cmd.Context()
			ref, err := parser.Parse(name)
			if err != nil {
				return err
			}

			local := local.NewController(conf.DefaultRoot())
			exist, err := local.Exist(ctx, ref, nil)
			if err != nil {
				return err
			}
			if exist {
				return errors.New("local repository already exists")
			}

			l := log.FromContext(ctx).WithFields(log.Fields{
				"ref": ref,
			})
			adaptor, ctrl, err := RemoteControllerFor(ctx, *tokens, ref)
			if err != nil {
				return fmt.Errorf("failed to get token for %s/%s: %w", ref.Host(), ref.Owner(), err)
			}

			// check repo has already existed
			if _, err := ctrl.Get(ctx, ref.Owner(), ref.Name(), nil); err == nil {
				l.Info("remote repository already exists")
			} else {
				me, err := ctrl.Me(ctx)
				if err != nil {
					return err
				}
				if f.Template == "" {
					ropt := &remote.CreateOption{
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
					}
					if me != ref.Owner() {
						ropt.Organization = ref.Owner()
					}

					if _, err := ctrl.Create(ctx, ref.Name(), ropt); err != nil {
						return err
					}
				} else {
					from, err := reporef.ParseSiblingRepoRef(ref, f.Template)
					if err != nil {
						return err
					}
					ropt := &remote.CreateFromTemplateOption{}
					if me != ref.Owner() {
						ropt.Owner = ref.Owner()
					}
					if f.Private {
						ropt.Private = true
					}
					if _, err = ctrl.CreateFromTemplate(ctx, from.Owner(), from.Name(), ref.Name(), ropt); err != nil {
						return err
					}
				}
			}
			accessToken, err := adaptor.GetAccessToken()
			if err != nil {
				l.WithField("error", err).Error("failed to get access token")
				return nil
			}
			for range f.CloneRetryLimit {
				_, err := local.Clone(ctx, ref, accessToken, nil)
				switch {
				case errors.Is(err, git.ErrRepositoryNotExists) || errors.Is(err, transport.ErrRepositoryNotFound):
					l.Info("waiting the remote repository is ready")
				case errors.Is(err, transport.ErrEmptyRemoteRepository):
					if _, err := local.Create(ctx, ref, nil); err != nil {
						l.WithField("error", err).
							WithField("error-type", fmt.Sprintf("%t", err)).
							Error("failed to create empty repository")
					} else {
						l.Info("created empty repository")
					}
					return nil
				case err == nil:
					return nil
				default:
					l.WithField("error", err).
						WithField("error-type", fmt.Sprintf("%t", err)).
						Error("failed to get repository")
					return nil
				}
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(1 * time.Second):
				}
			}
			return nil
		},
	}
	defaults.Create.CloneRetryLimit = 5

	cmd.Flags().
		BoolVarP(&f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().
		StringVarP(&f.Template, "template", "", defaults.Create.Template, "Create new repository from the template")
	cmd.Flags().
		StringVarP(&f.Description, "description", "", "", "A short description of the repository")
	cmd.Flags().
		StringVarP(&f.Homepage, "homepage", "", "", "A URL with more information about the repository")
	cmd.Flags().
		StringVarP(&f.LicenseTemplate, "license-template", "", defaults.Create.LicenseTemplate, `Choose an open source license template that best suits your needs, and then use the license keyword as the license_template string when "auto-init" flag is set. For example, "mit" or "mpl-2.0"`)
	cmd.Flags().
		StringVarP(&f.GitignoreTemplate, "gitignore-template", "", defaults.Create.GitignoreTemplate, `Desired language or platform .gitignore template to apply when "auto-init" flag is set. Use the name of the template without the extension. For example, "Haskell"`)
	cmd.Flags().
		BoolVarP(&f.Private, "private", "", defaults.Create.Private, "Whether the repository is private")
	cmd.Flags().
		BoolVarP(&f.IsTemplate, "is-template", "", false, "Whether the repository is available as a template")
	cmd.Flags().
		BoolVarP(&f.DisableDownloads, "disable-downloads", "", defaults.Create.DisableDownloads, `Disable "Downloads" page`)
	cmd.Flags().
		BoolVarP(&f.DisableWiki, "disable-wiki", "", defaults.Create.DisableWiki, `Disable Wiki for the repository`)
	cmd.Flags().
		BoolVarP(&f.AutoInit, "auto-init", "", defaults.Create.AutoInit, "Create an initial commit with empty README")
	cmd.Flags().
		BoolVarP(&f.DisableProjects, "disable-projects", "", defaults.Create.DisableProjects, `Disable projects for the repository`)
	cmd.Flags().
		BoolVarP(&f.DisableIssues, "disable-issues", "", defaults.Create.DisableIssues, `Disable issues for the repository`)
	cmd.Flags().
		BoolVarP(&f.PreventSquashMerge, "prevent-squash-merge", "", defaults.Create.PreventSquashMerge, "Prevent squash-merging pull requests")
	cmd.Flags().
		BoolVarP(&f.PreventMergeCommit, "prevent-merge-commit", "", defaults.Create.PreventMergeCommit, "Prevent merging pull requests with a merge commit")
	cmd.Flags().
		BoolVarP(&f.PreventRebaseMerge, "prevent-rebase-merge", "", defaults.Create.PreventRebaseMerge, "Prevent rebase-merging pull requests")
	cmd.Flags().
		BoolVarP(&f.DeleteBranchOnMerge, "delete-branch-on-merge", "", defaults.Create.DeleteBranchOnMerge, "Allow automatically deleting head branches when pull requests are merged")
	cmd.Flags().
		IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", defaults.Create.CloneRetryLimit, "")
	return cmd
}
