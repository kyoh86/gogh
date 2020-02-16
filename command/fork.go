package command

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gogh/gogh"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx gogh.Context, gitClient GitClient, hubClient HubClient, update, withSSH, shallow bool, organization string, repo *gogh.Repo) error {
	log.Print("info: Finding a repository")
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}
	log.Print("info: Getting a repository")
	if !project.Exists {
		repoURL := repo.URL(withSSH)
		log.Print("info: Clone", fmt.Sprintf("%s -> %s", repoURL, project.FullPath))
		if err := gitClient.Clone(project.FullPath, repoURL, shallow); err != nil {
			return err
		}
	} else if update {
		log.Print("info: Update", project.FullPath)
		if err := gitClient.Update(project.FullPath); err != nil {
			return err
		}
	}
	log.Print("info: Forking a repository")
	newRepo, err := hubClient.Fork(ctx, repo, organization)
	if err != nil {
		var accepted *github.AcceptedError
		if !errors.As(err, &accepted) {
			return err
		}
	}
	log.Print("info: Getting remotes")
	remotes, err := gitClient.GetRemotes(project.FullPath)
	if err != nil {
		return err
	}

	log.Print("info: Removing old remotes")
	owner := repo.Owner()
	me := newRepo.Owner()
	for name := range remotes {
		if name == me || name == owner || name == "origin" {
			if err := gitClient.RemoveRemote(project.FullPath, name); err != nil {
				return err
			}
		}
	}

	log.Print("info: Creating new remotes")
	if err := gitClient.AddRemote(project.FullPath, owner, repo.URL(withSSH)); err != nil {
		return err
	}
	if err := gitClient.AddRemote(project.FullPath, me, newRepo.URL(withSSH)); err != nil {
		return err
	}

	log.Print("info: Fetching new remotes")
	if err := gitClient.Fetch(project.FullPath); err != nil {
		return err
	}

	log.Printf("info: Setting upstream to %q", me)
	branch, err := gitClient.GetCurrentBranch(project.FullPath)
	if err != nil {
		return err
	}
	if err := gitClient.SetUpstreamTo(project.FullPath, me+"/"+branch); err != nil {
		return err
	}
	return nil
}
