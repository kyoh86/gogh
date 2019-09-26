package command

import (
	"fmt"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/remote"
)

// Repos will show a list of repositories for a user.
func Repos(ctx gogh.Context, user string, own, collaborate, member bool, visibility, sort, direction string) error {
	repos, err := remote.Repos(ctx, user, own, collaborate, member, visibility, sort, direction)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Println(repo)
	}
	return nil
}