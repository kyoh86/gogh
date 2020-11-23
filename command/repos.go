package command

import (
	"context"
	"os"
	"strings"

	"github.com/kyoh86/gogh/gogh"
)

// Repos will show a list of repositories for a user.
func Repos(
	ctx context.Context,
	ev gogh.Env,
	hubClient HubClient,
	user string,
	own,
	collaborate,
	member,
	archived bool,
	visibility,
	sort,
	direction string,
	formatter gogh.ProjectListFormatter,
	query string,
) error {
	repos, err := hubClient.Repos(ctx, ev, user, own, collaborate, member, archived, visibility, sort, direction)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		spec, err := gogh.ParseRepoSpec(repo)
		if err != nil {
			continue
		}
		p, err := gogh.NewProject(ev, spec)
		if err != nil {
			continue
		}
		if query != "" && query != p.FullPath && !strings.Contains(p.RelPath, query) {
			continue
		}
		formatter.Add(p)
	}
	return formatter.PrintAll(os.Stdout, "\n")
}
