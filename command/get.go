package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx gogh.Context, update, withSSH, shallow bool, repos gogh.Repos) error {
	for _, repo := range repos {
		repo := repo
		if err := Get(ctx, update, withSSH, shallow, &repo); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(ctx gogh.Context, update, withSSH, shallow bool, repo *gogh.Repo) error {
	repoURL := repo.URL(ctx, withSSH)
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}
	if !project.Exists {
		log.Println("info: Clone", fmt.Sprintf("%s -> %s", repoURL, project.FullPath))
		return gitClone(ctx, repoURL, project.FullPath, shallow)
	}
	if update {
		log.Println("info: Update", project.FullPath)
		return gitUpdate(ctx, project.FullPath)
	}
	log.Println("warn: Exists", project.FullPath)
	return nil
}
