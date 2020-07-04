package command

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/gogh"
)

// Repos will show a list of repositories for a user.
func Repos(ctx context.Context, ev gogh.Env, hubClient HubClient, user string, own, collaborate, member, archived bool, visibility, sort, direction string) error {
	repos, err := hubClient.Repos(ctx, ev, user, own, collaborate, member, archived, visibility, sort, direction)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Println(repo)
	}
	return nil
}
