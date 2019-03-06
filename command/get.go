package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx gogh.Context, update, withSSH, shallow bool, remotes gogh.Remotes) error {
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
func Get(ctx gogh.Context, update, withSSH, shallow bool, remote *gogh.Remote) error {
	remoteURL := remote.URL(ctx, withSSH)
	project, err := gogh.FindProject(ctx, remote)
	if err != nil {
		return err
	}

	if !project.Exists {
		log.Println("info: clone", fmt.Sprintf("%s -> %s", remoteURL, project.FullPath))
		return gitClone(ctx, remoteURL, project.FullPath, shallow)
	}
	if update {
		log.Println("info: update", project.FullPath)
		return gitUpdate(ctx, project.FullPath)
	}
	log.Println("warn: exists", project.FullPath)
	return nil
}
