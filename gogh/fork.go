package gogh

import (
	"fmt"
	"log"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx Context, update, withSSH, shallow, noRemote bool, remoteName string, organization string, remote *Remote) error {
	log.Printf("info: cloning a repository")
	if err := Get(ctx, update, withSSH, shallow, remote); err != nil {
		return err
	}

	local, err := FindLocal(ctx, remote)
	if err != nil {
		return err
	}
	log.Printf("info: forking a repository")
	if err := hubFork(ctx, local, remote, noRemote, remoteName, organization); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(ctx.Stdout(), local.RelPath); err != nil {
		return err
	}
	return nil
}
