package command

import (
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// Repos will show a list of repositories for a user.
func Repos(ctx gogh.Context, hubClient HubClient, user string, own, collaborate, member bool, visibility, sort, direction string) error {
	InitLog(ctx)

	repos, err := hubClient.Repos(ctx, user, own, collaborate, member, visibility, sort, direction)
	if err != nil {
		return err
	}
	var stdout io.Writer = os.Stdout
	if ctx, ok := ctx.(gogh.IOContext); ok {
		stdout = ctx.Stdout()
	}
	for _, repo := range repos {
		fmt.Fprintln(stdout, repo)
	}
	return nil
}
