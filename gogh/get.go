package gogh

import (
	"fmt"
	"log"
	"os"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx Context, update, withSSH, shallow bool, remoteNames RemoteNames) error {
	for _, remoteName := range remoteNames {
		if err := Get(ctx, update, withSSH, shallow, remoteName); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(ctx Context, update, withSSH, shallow bool, remoteName RemoteName) error {
	remoteURL := remoteName.URL(ctx, withSSH)
	repo, err := FromURL(ctx, remoteURL)
	if err != nil {
		return err
	}

	path := repo.FullPath
	newPath := false
	_, err = os.Stat(path)
	switch {
	case err == nil:
		// noop
	case os.IsNotExist(err):
		newPath = true
	default:
		return err
	}

	if newPath {
		log.Println("info: clone", fmt.Sprintf("%s -> %s", remoteURL, path))

		return gitClone(ctx, remoteURL, path, shallow)
	}
	if update {
		log.Println("info: update", path)
		return gitUpdate(ctx, path)
	}
	log.Println("warn: exists", path)
	return nil
}
