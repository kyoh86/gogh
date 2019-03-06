package remote

import (
	"fmt"

	"github.com/kyoh86/gogh/gogh"
	remo "github.com/kyoh86/gogh/gogh/remote"
)

// Repo will show a list of repositories for a user.
func Repo(ctx gogh.Context, user string, own, collaborate, member bool, visibility, sort, direction string) error {
	repos, err := remo.Repo(ctx, user, own, collaborate, member, visibility, sort, direction)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Println(repo)
	}
	return nil
}
