package gogh

import (
	"fmt"
	"log"
	"os"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx Context, update, withSSH, shallow bool, repoSpecs RepoSpecs) error {
	for _, repoSpec := range repoSpecs {
		if err := Get(ctx, update, withSSH, shallow, repoSpec); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(ctx Context, update, withSSH, shallow bool, repoSpec RepoSpec) error {
	remoteURL := repoSpec.URL(ctx, withSSH)
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

		return gitClone(remoteURL, path, shallow)
	}
	if update {
		log.Println("info: update", path)
		return gitUpdate(path)
	}
	log.Println("warn: exists", path)
	return nil
}
