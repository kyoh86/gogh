package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var deleteFlags struct {
	force bool
}

var deleteCommand = &cobra.Command{
	Use:     "delete [flags] [[OWNER/]NAME]",
	Aliases: []string{"remove"},
	Short:   "Delete a repository with a remote repository",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		ctx := cmd.Context()
		servers := app.Servers()
		var selected string
		if len(specs) == 0 {
			servers, err := servers.List()
			if err != nil {
				return err
			}
			for _, server := range servers {
				adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
				if err != nil {
					return err
				}
				remote := gogh.NewRemoteController(adaptor)
				founds, err := remote.List(ctx, nil)
				if err != nil {
					return err
				}
				for _, s := range founds {
					specs = append(specs, s.String())
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

		parser := gogh.NewSpecParser(servers)
		spec, server, err := parser.Parse(selected)
		if err != nil {
			return err
		}

		local := gogh.NewLocalController(app.DefaultRoot())
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
		if err := local.Delete(ctx, spec, nil); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("delete local: %w", err)
			}
		}

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
		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		if err := gogh.NewRemoteController(adaptor).Delete(ctx, spec.Owner(), spec.Name(), nil); err != nil {
			var gherr *github.ErrorResponse
			if errors.As(err, &gherr) && gherr.Response.StatusCode == http.StatusForbidden {
				log.FromContext(ctx).Errorf("Failed to delete a repository: there is no permission to delete %q", spec.URL())
				log.FromContext(ctx).Errorf(`Add scope "delete_repo" for the token for %q`, server.String())
			} else {
				return err
			}
		}
		return nil
	},
}

func init() {
	deleteCommand.Flags().BoolVarP(&deleteFlags.force, "force", "", false, "Do NOT confirm to delete.")
	facadeCommand.AddCommand(deleteCommand)
}
