package command

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/remote"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx gogh.Context, update, withSSH, shallow bool, organization string, repo *gogh.Repo) error {
	log.Printf("info: Finding a repository")
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}
	log.Printf("info: Getting a repository")
	if !project.Exists {
		repoURL := repo.URL(ctx, withSSH)
		log.Println("info: Clone", fmt.Sprintf("%s -> %s", repoURL, project.FullPath))
		if err := git().Clone(ctx, project, repoURL, shallow); err != nil {
			return err
		}
	} else if update {
		log.Println("info: Update", project.FullPath)
		if err := git().Update(ctx, project); err != nil {
			return err
		}
	}
	log.Printf("info: Forking a repository")
	newRepo, err := remote.Fork(ctx, repo, organization)
	if err != nil {
		var accepted *github.AcceptedError
		if !errors.As(err, &accepted) {
			return err
		}
	}
	remotes, err := git().GetRemotes(ctx, project)
	if err != nil {
		return err
	}
	for name := range remotes {
		if name == newRepo.Owner(ctx) || name == repo.Owner(ctx) || name == "origin" {
			if err := git().RemoveRemote(ctx, project, name); err != nil {
				return err
			}
		}
	}
	if err := git().AddRemote(ctx, project, repo.Owner(ctx), repo.URL(ctx, withSSH)); err != nil {
		return err
	}
	if err := git().AddRemote(ctx, project, newRepo.Owner(ctx), newRepo.URL(ctx, withSSH)); err != nil {
		return err
	}
	fmt.Fprintln(ctx.Stdout(), project.FullPath)
	return nil
}
