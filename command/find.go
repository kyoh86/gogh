package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Find a path of a project
func Find(ctx gogh.Context, repo *gogh.Repo) error {
	log.Println("info: Finding a repository...")
	path, err := gogh.FindProjectPath(ctx, repo)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(ctx.Stdout(), path); err != nil {
		return err
	}

	return nil
}
