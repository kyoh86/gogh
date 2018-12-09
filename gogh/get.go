package gogh

import (
	"fmt"
	"log"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx Context, update, withSSH, shallow bool, remotes Remotes) error {
	for _, remote := range remotes {
		remote := remote
		if err := Get(ctx, update, withSSH, shallow, &remote); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(ctx Context, update, withSSH, shallow bool, remote *Remote) error {
	remoteURL := remote.URL(ctx, withSSH)
	local, err := FindLocal(ctx, remote)
	if err != nil {
		return err
	}

	if !local.Exists {
		log.Println("info: clone", fmt.Sprintf("%s -> %s", remoteURL, local.FullPath))
		return gitClone(ctx, remoteURL, local.FullPath, shallow)
	}
	if update {
		log.Println("info: update", local.FullPath)
		return gitUpdate(ctx, local.FullPath)
	}
	log.Println("warn: exists", local.FullPath)
	return nil
}
