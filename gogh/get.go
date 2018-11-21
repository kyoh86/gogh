package gogh

import (
	"fmt"
	"log"
	"os"

	"github.com/kyoh86/gogh/internal/git"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx Context, update, withSSH, shallow bool, repoSpecs Specs) error {
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
		log.Println("clone", fmt.Sprintf("%s -> %s", remoteURL, path))

		return git.Clone(remoteURL, path, shallow)
	}
	if update {
		log.Println("update", path)
		return git.Update(path)
	}
	log.Println("exists", path)
	return nil
}
