package pr

import (
	"os"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/kyoh86/gogh/cl"
	"github.com/kyoh86/gogh/gh"
	"github.com/kyoh86/gogh/gh/flags"
	"github.com/kyoh86/gogh/util"
	"gopkg.in/alecthomas/kingpin.v2"
)

// List pull requests
func List(c *kingpin.CmdClause) gh.Command {
	var (
		owner string
		repos string

		ops github.PullRequestListOptions

		rowFormat string
	)

	flags.Repository(c, &owner, &repos)

	flags.Sort(c).EnumVar(&ops.Sort, "closed", "created", "updated", "popularity", "long-running")
	flags.Direction(c).EnumVar(&ops.Direction, "asc", "desc")
	flags.PerPage(c).IntVar(&ops.PerPage)
	flags.Page(c).IntVar(&ops.PerPage)

	c.Flag("state", "Either open, closed, or all to filter by state").Default("all").EnumVar(&ops.State, "open", "closed", "all")
	c.Flag("head", "Filter pulls by head user and branch name in the format of user:ref-name").StringVar(&ops.Head)
	c.Flag("base", "Filter pulls by base branch name").StringVar(&ops.Base)

	c.Flag("row-format", "Format to output with [gtf](https://github.com/leekchan/gtf)").Default(strings.Join([]string{
		`#{{.Number}}`,
		`{{.Title}}`,
		`{{.Base.Ref}}`,
		`{{.CreatedAt | date "01-02 15:04"}}`,
		`{{.MergedAt | date "01-02 15:04"}}`,
		`{{.ClosedAt | date "01-02 15:04"}}`,
	}, "\t") + "\n").StringVar(&rowFormat)

	return func() error {
		logrus.Debugf("running on %s/%s", owner, repos)

		t := flags.Template()
		formatter, err := t.Parse(rowFormat)
		if err != nil {
			return util.WrapErr("Failed to parse row-format as template", err)
		}

		client, err := cl.GitHubClient()
		if err != nil {
			return err
		}

		order := ops.Sort
		if order == "closed" {
			ops.Sort = ""
		}

		requests, _, err := client.PullRequests.List(owner, repos, &ops)
		if err != nil {
			return util.WrapErr("Failed to list up pulls", err)
		}

		pulls := &list{order: order, direction: ops.Direction, array: requests}
		sort.Sort(pulls)

		for _, request := range pulls.array {
			formatter.Execute(os.Stdout, request)
		}
		return nil
	}
}

type list struct {
	order     string
	direction string
	array     []github.PullRequest
}

func (p *list) Len() int {
	return len(p.array)
}

func (p *list) Swap(i, j int) {
	p.array[i], p.array[j] = p.array[j], p.array[i]
}

func (p *list) Less(i, j int) bool {
	switch p.order {
	case "closed":
		if p.array[i].ClosedAt == nil {
			return p.array[j].ClosedAt != nil
		}
		if p.array[j].ClosedAt == nil {
			return false
		}
		if p.direction == "desc" {
			return p.array[i].ClosedAt.After(*p.array[j].ClosedAt)
		}
		return p.array[i].ClosedAt.Before(*p.array[j].ClosedAt)
	default:
		return false
	}
}
