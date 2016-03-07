package pr

import (
	"fmt"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/kyoh86/gogh/internal/cl"
	"github.com/kyoh86/gogh/internal/gh"
	"github.com/kyoh86/gogh/internal/util"
	"gopkg.in/alecthomas/kingpin.v2"
)

// CreateCommand creates pull request
func CreateCommand(c *kingpin.CmdClause) gh.Command {
	var (
		ops = github.NewPullRequest{
			Title: util.StringPtr(""),
			Body:  util.StringPtr(""),
		}
		head gh.Branch
		base BaseBranch
	)

	headFlag := c.Flag("head", "The name of the repository/branch where your changes are implemented").Short('h')
	workingBranch, working := gh.WorkingBranch()
	if working {
		headFlag.Default(workingBranch.String())
	} else {
		headFlag.Required()
	}
	headFlag.SetValue(&head)

	baseFlag := c.Flag("base", "The name of the branch you want your changes pulled into").Short('b').Required()
	baseFlag.SetValue(&base)

	c.Flag("title", "The title of the pull request").Short('m').Required().StringVar(ops.Title)
	c.Flag("body", "The contents of the pull request").StringVar(ops.Body)

	return func() error {
		logrus.Debugf("running on %s", head.String())

		owner := base.Owner
		if owner == "" {
			owner = head.Owner
			ops.Head = &head.Branch
		} else {
			ops.Head = util.StringPtr(head.Owner + ":" + head.Branch)
		}
		ops.Base = util.StringPtr(base.String())

		client, err := cl.GitHubClient()
		if err != nil {
			return err
		}

		pr, _, err := client.PullRequests.Create(owner, head.Repo, &ops)
		if err != nil {
			return err
		}

		fmt.Println(*pr.HTMLURL)
		return nil
	}
}

var (
	// BaseBranchRegexp for text intending a base-branch
	BaseBranchRegexp = regexp.MustCompile(`^(?:(?P<owner>` + gh.NamePattern + `):)?(?P<branch>` + gh.NamePattern + `)$`)
)

// BaseBranch contains identifier of pull-request target branch: "[owner:]branch"
type BaseBranch struct {
	Owner  string
	Branch string
}

// Set a value string to BaseBranch
func (r *BaseBranch) Set(value string) error {
	names := BaseBranchRegexp.SubexpNames()
	match := BaseBranchRegexp.FindStringSubmatch(value)
	if len(match) < len(names) {
		return fmt.Errorf("specified parameter '%s' is not a branch", value)
	}
	for i, name := range names {
		if match[i] == "" {
			continue
		}
		switch name {
		case "owner":
			r.Owner = match[i]
		case "branch":
			r.Branch = match[i]
		}
	}
	return nil
}

func (r *BaseBranch) String() string {
	if r.Owner == "" {
		return r.Branch
	}
	return fmt.Sprintf("%s:%s", r.Owner, r.Branch)
}
