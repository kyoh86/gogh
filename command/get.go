package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// GetAll clonse or updates remote repositories.
func GetAll(ev gogh.Env, gitClient GitClient, update, withSSH, shallow bool, specs gogh.RepoSpecs) error {
	for _, spec := range specs {
		spec := spec
		if err := Get(ev, gitClient, update, withSSH, shallow, &spec); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(ev gogh.Env, gitClient GitClient, update, withSSH, shallow bool, spec *gogh.RepoSpec) error {
	project, repo, err := gogh.FindOrNewProject(ev, spec)
	if err != nil {
		return err
	}
	if !project.Exists {
		repoURL := repo.URL(withSSH)
		log.Println("info: Clone", fmt.Sprintf("%s -> %s", repoURL, project.FullPath))
		if err := gitClient.Clone(project.FullPath, repoURL, shallow); err != nil {
			return err
		}
		fmt.Println(project.FullPath)
		return nil
	}
	if update {
		log.Println("info: Update", project.FullPath)
		if err := gitClient.Update(project.FullPath); err != nil {
			return err
		}
		fmt.Println(project.FullPath)
	}
	log.Println("warn: Exists", project.FullPath)
	return nil
}
