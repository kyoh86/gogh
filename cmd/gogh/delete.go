package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var deleteFlags struct {
	local  bool
	remote bool
	force  bool
	dryrun bool
}

var deleteCommand = &cobra.Command{
	Use:     "delete [flags] [[OWNER/]NAME]",
	Aliases: []string{"remove"},
	Short:   "Delete a project with a remote repository",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		ctx := cmd.Context()
		var selected string
		if len(specs) == 0 {
			for _, entry := range tokens.Entries() {
				adaptor, err := github.NewAdaptor(ctx, entry.Host, string(entry.Token))
				if err != nil {
					return err
				}
				remote := gogh.NewRemoteController(adaptor)
				founds, err := remote.List(ctx, nil)
				if err != nil {
					return err
				}
				for _, s := range founds {
					specs = append(specs, s.Spec.String())
				}
			}
			if err := survey.AskOne(&survey.Select{
				Message: "A repository to delete",
				Options: specs,
			}, &selected); err != nil {
				return err
			}
		} else {
			selected = specs[0]
		}

		parser := gogh.NewSpecParser(tokens.GetDefaultKey())
		spec, err := parser.Parse(selected)
		if err != nil {
			return err
		}

		if deleteFlags.local {
			local := gogh.NewLocalController(defaultRoot())
			if !deleteFlags.force {
				var confirmed bool
				if err := survey.AskOne(&survey.Confirm{
					Message: fmt.Sprintf("Are you sure you want to delete local-project %s?", spec.String()),
				}, &confirmed); err != nil {
					return err
				}
				if !confirmed {
					return nil
				}
			}
			if deleteFlags.dryrun {
				fmt.Printf("deleting local %s\n", spec.String())
			} else if err := local.Delete(ctx, spec, nil); err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("delete local: %w", err)
				}
			}
		}

		if deleteFlags.remote {
			if !deleteFlags.force {
				var confirmed bool
				if err := survey.AskOne(&survey.Confirm{
					Message: fmt.Sprintf("Are you sure you want to delete remote-repository %s?", spec.String()),
				}, &confirmed); err != nil {
					return err
				}
				if !confirmed {
					return nil
				}
			}
			token := tokens.Get(spec.Host(), spec.Owner())
			adaptor, err := github.NewAdaptor(ctx, spec.Host(), string(token))
			if err != nil {
				return err
			}
			if deleteFlags.dryrun {
				fmt.Printf("deleting remote %s\n", spec.String())
			} else if err := gogh.NewRemoteController(adaptor).Delete(ctx, spec.Owner(), spec.Name(), nil); err != nil {
				var gherr *github.ErrorResponse
				if errors.As(err, &gherr) && gherr.Response.StatusCode == http.StatusForbidden {
					log.FromContext(ctx).Errorf("Failed to delete a repository: there is no permission to delete %q", spec.URL())
					log.FromContext(ctx).Error(`Add scope "delete_repo" for the token`)
				} else {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	deleteCommand.Flags().BoolVarP(&deleteFlags.local, "local", "", true, "Delete local project.")
	deleteCommand.Flags().
		BoolVarP(&deleteFlags.remote, "remote", "", false, "Delete remote project.")
	deleteCommand.Flags().
		BoolVarP(&deleteFlags.force, "force", "", false, "Do NOT confirm to delete.")
	deleteCommand.Flags().
		BoolVarP(&deleteFlags.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	facadeCommand.AddCommand(deleteCommand)
}
