package command

import (
	"fmt"

	"github.com/kyoh86/gogh/gogh"
)

// Repos will show a list of repositories for a user.
func Repos(ctx gogh.Context, hubClient HubClient, user string, own, collaborate, member bool, visibility, sort, direction string) error {
	if err := InitLog(ctx); err != nil {
		return err
	}

	repos, err := hubClient.Repos(ctx, user, own, collaborate, member, visibility, sort, direction)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Fprintln(ctx.Stdout(), repo)
	}
	return nil
}
