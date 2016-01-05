package pr

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/kyoh86/gogh/cl"
	"github.com/kyoh86/gogh/gh"
	"github.com/kyoh86/gogh/gh/flags"
	"github.com/kyoh86/gogh/util"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Create pull requests
func Create(c *kingpin.CmdClause) gh.Command {
	var (
		owner string
		repos string

		ops github.NewPullRequest
	)

	ops = github.NewPullRequest{
		Title: util.StringPtr(""),
		Head:  util.StringPtr(""),
		Base:  util.StringPtr(""),
		Body:  util.StringPtr(""),
	}

	flags.Repository(c, &owner, &repos)
	flags.HeadBranch(c).StringVar(ops.Head)

	c.Flag("title", "The title of the pull request").Short('m').Required().StringVar(ops.Title)
	c.Flag("base", "The name of the branch you want your changes pulled into").Short('b').Required().StringVar(ops.Base)
	c.Flag("body", "The contents of the pull request").StringVar(ops.Body)

	return func() error {
		logrus.Debugf("running on %s/%s", owner, repos)

		client, err := cl.GitHubClient()
		if err != nil {
			return err
		}

		pr, _, err := client.PullRequests.Create(owner, repos, &ops)
		if err != nil {
			return err
		}

		fmt.Println(*pr.HTMLURL)
		return nil
	}
}
