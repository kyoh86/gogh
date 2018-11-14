package gogh

import (
	"fmt"
	"log"
	"os"

	"github.com/kyoh86/gogh/internal/git"
	"github.com/kyoh86/gogh/repo"
)

// GetAll clonse or updates remote repositories.
func GetAll(update, withSSH, shallow bool, repoSpecs repo.Specs) error {
	for _, repoSpec := range repoSpecs {
		if err := Get(update, withSSH, shallow, repoSpec); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(update, withSSH, shallow bool, repoSpec repo.Spec) error {
	rmt, err := repoSpec.Remote(withSSH)
	if err != nil {
		return err
	}

	remoteURL := rmt.URL()
	local, err := repo.FromURL(remoteURL)
	if err != nil {
		return err
	}

	path := local.FullPath
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
