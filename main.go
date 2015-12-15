package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/github"
	"github.com/kyoh86/ghu/gh"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Conf struct {
	Token      string `json:"token"`
	Owner      string `json:"owner"`
	Repository string `json:"repository"`

	Methods struct {
		PullRequest struct {
			List struct {
				github.PullRequestListOptions
				Limit  int
				Header bool
			} `json:"list"`
		} `json:"pullRequest"`
	} `json:"methods"`
}

func main() {

	conf := Conf{}
	app := kingpin.New("kinpin.App", "Sample for kingpin.App")
	app.Flag("access-token", "Access token for GitHub").Short('t').Required().StringVar(&conf.Token)

	pr := app.Command("pull-request", "Access to pull-requests")
	prlist := pr.Command("list", "List up pull-requests")
	prlist.Flag("owner", "Repository owner name").Short('o').Required().StringVar(&conf.Owner)
	prlist.Flag("repos", "Repository name").Short('r').Required().StringVar(&conf.Repository)
	prlist.Flag("state", "State of pull-request").EnumVar(&conf.Methods.PullRequest.List.State, "open", "closed")
	prlist.Flag("head", "Head of pull-request").StringVar(&conf.Methods.PullRequest.List.Head)
	prlist.Flag("base", "Base branch name of pull-request").StringVar(&conf.Methods.PullRequest.List.Base)
	prlist.Flag("sort", "Sort specifies how to sort pull requests").EnumVar(&conf.Methods.PullRequest.List.Sort, "closed", "created", "updated", "popularity", "long-running")
	prlist.Flag("direction", "Direction in which to sort pull requests").EnumVar(&conf.Methods.PullRequest.List.Direction, "asc", "desc")

	prlist.Flag("header", "Repository owner name").Default("true").BoolVar(&conf.Methods.PullRequest.List.Header)
	prlist.Flag("limit", "Limit to get a pull-requests").IntVar(&conf.Methods.PullRequest.List.Limit)
	// TODO: access-tokenを保存できるようにする: https://github.com/tcnksm/go-gitconfig とか使えば容易かも

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case prlist.FullCommand():
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: conf.Token},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client := github.NewClient(tc)

		// list all pull-requests for the authenticated user and owner/repository
		order := conf.Methods.PullRequest.List.PullRequestListOptions.Sort
		if order == "closed" {
			conf.Methods.PullRequest.List.PullRequestListOptions.Sort = ""
		}
		requests, _, err := client.PullRequests.List(conf.Owner, conf.Repository, &conf.Methods.PullRequest.List.PullRequestListOptions)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		pulls := gh.NewPullRequests(requests)
		pulls.Order = gh.PullRequestSortProperty(order)
		sort.Sort(pulls)

		if conf.Methods.PullRequest.List.Header {
			fmt.Println("#\ttitle\tbase\topen\tmerge\tclose")
		}
		now := time.Now()
		in := func(t *time.Time) string {
			if t == nil {
				return "-"
			}
			return t.In(now.Location()).Format("01-02 15:04")
		}

		// TODO :表示する項目の選択をできるようにする
		for i, request := range pulls.Array {
			if i > conf.Methods.PullRequest.List.Limit && conf.Methods.PullRequest.List.Limit > 0 {
				break
			}
			fmt.Printf("#%d\t%s\t%s\t%s\t%s\t%s\n", *request.Number, *request.Title, *request.Base.Ref, in(request.CreatedAt), in(request.MergedAt), in(request.ClosedAt))
		}
	}
}
