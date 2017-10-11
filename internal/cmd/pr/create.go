package pr

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/kyoh86/gogh/internal/cmd"
	"github.com/kyoh86/gogh/internal/flags"
	"github.com/kyoh86/gogh/internal/repo"
	"github.com/sirupsen/logrus"
	"github.com/wacul/ptr"
	"gopkg.in/alecthomas/kingpin.v2"
)

// CreateCommand creates pull request
func CreateCommand(c *kingpin.CmdClause, r *repo.Repository) cmd.Command {
	var (
		ops = github.NewPullRequest{
			Title: ptr.String(""),
			Body:  ptr.String(""),
		}
		head flags.Name
		base flags.BaseBranch
	)

	headFlag := c.Flag("head", "The name of the repository/branch where your changes are implemented").Short('h')
	branch, err := r.Branch()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get branch")
	}
	headFlag.Default(branch)
	headFlag.SetValue(&head)

	baseFlag := c.Flag("base", "The name of the branch you want your changes pulled into").Short('b').Default("master")
	baseFlag.SetValue(&base)

	c.Flag("title", "The title of the pull request").Short('m').StringVar(ops.Title)
	c.Flag("body", "The contents of the pull request").StringVar(ops.Body)

	return func(ctx context.Context) error {
		id, err := r.Identifier()
		if err != nil {
			return err
		}
		ops.Head = ptr.String(head.String())
		ops.Base = ptr.String(base.String())

		client, err := cmd.GitHubClient()
		if err != nil {
			return err
		}

		pr, _, err := client.PullRequests.Create(ctx, id.Owner, id.Name, &ops)
		if err != nil {
			return err
		}

		fmt.Println(*pr.HTMLURL)
		return nil
	}
}
