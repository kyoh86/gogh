package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx gogh.Context, update, withSSH, shallow, noRemote bool, remoteName string, organization string, remote *gogh.Remote) error {
	log.Printf("info: cloning a repository")
	if err := Get(ctx, update, withSSH, shallow, remote); err != nil {
		return err
	}

	project, err := gogh.FindProject(ctx, remote)
	if err != nil {
		return err
	}
	log.Printf("info: forking a repository")
	if err := hubFork(ctx, project, remote, noRemote, remoteName, organization); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(ctx.Stdout(), project.RelPath); err != nil {
		return err
	}
	return nil
}
