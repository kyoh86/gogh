package command

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// GetAll clonse or updates remote repositories.
func GetAll(ctx gogh.Context, gitClient GitClient, update, withSSH, shallow bool, repos gogh.Repos) error {
	InitLog(ctx)

	for _, repo := range repos {
		repo := repo
		if err := Get(ctx, gitClient, update, withSSH, shallow, &repo); err != nil {
			return err
		}
	}
	return nil
}

// Get clones or updates a remote repository.
// If update is true, updates the locally cloned repository. Otherwise does nothing.
// If shallow is true, does shallow cloning. (no effect if already cloned or the VCS is Mercurial and git-svn)
func Get(ctx gogh.Context, gitClient GitClient, update, withSSH, shallow bool, repo *gogh.Repo) error {
	InitLog(ctx)

	repoURL := repo.URL(ctx, withSSH)
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}
	var stdout io.Writer = os.Stdout
	if ctx, ok := ctx.(gogh.IOContext); ok {
		stdout = ctx.Stdout()
	}
	if !project.Exists {
		log.Println("info: Clone", fmt.Sprintf("%s -> %s", repoURL, project.FullPath))
		if err := gitClient.Clone(project.FullPath, repoURL, shallow); err != nil {
			return err
		}
		fmt.Fprintln(stdout, project.FullPath)
		return nil
	}
	if update {
		log.Println("info: Update", project.FullPath)
		if err := gitClient.Update(project.FullPath); err != nil {
			return err
		}
		fmt.Fprintln(stdout, project.FullPath)
	}
	log.Println("warn: Exists", project.FullPath)
	return nil
}
