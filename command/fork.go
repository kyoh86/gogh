package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx gogh.Context, update, withSSH, shallow, noRemote bool, remoteName string, organization string, repo *gogh.Repo) error {
	log.Printf("info: Finding a repository")
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}
	log.Printf("info: Forking a repository")
	if err := hub().Fork(ctx, project, repo, noRemote, remoteName, organization); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(ctx.Stdout(), project.RelPath); err != nil {
		return err
	}
	return nil
}
