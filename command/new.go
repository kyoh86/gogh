package command

import (
	"log"
	"net/url"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// New creates a local project and a remote repository.
func New(
	ctx gogh.Context,
	gitClient GitClient,
	hubClient HubClient,
	private bool,
	description string,
	homepage *url.URL,
	bare bool,
	template string,
	separateGitDir string,
	shared RepoShared,
	repo *gogh.Repo,
) error {
	log.Printf("info: Creating new project and a remote repository %s", repo)
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}

	log.Printf("info: Checking existing project")
	remote, err := checkProjectRemote(ctx, gitClient, project, repo)
	if err != nil {
		return err
	}

	// mkdir
	log.Println("info: Creating a directory")
	if err := os.MkdirAll(project.FullPath, os.ModePerm); err != nil {
		return err
	}

	// git init
	log.Println("info: Initializing a repository")
	if err := gitClient.Init(project.FullPath, bare, template, separateGitDir, shared.String()); err != nil {
		return err
	}

	if remote {
		return nil
	}

	// hub create
	log.Println("info: Creating a new repository in GitHub")
	newRepo, err := hubClient.Create(ctx, repo, description, homepage, private)
	if err != nil {
		return err
	}

	// git remote add origin
	url, err := url.Parse(newRepo.GetHTMLURL())
	if err != nil {
		return err
	}
	if err := gitClient.AddRemote(project.FullPath, "origin", url); err != nil {
		return err
	}

	return nil
}

func checkProjectRemote(ctx gogh.Context, gitClient GitClient, project *gogh.Project, repo *gogh.Repo) (bool, error) {
	if !project.Exists {
		return false, nil
	}
	remotes, err := gitClient.GetRemotes(project.FullPath)
	if err != nil {
		return false, err
	}
	if len(remotes) > 0 {
		remote := remotes["origin"]
		if remote == nil {
			return false, nil
		}
		if remote.String() == repo.URL(ctx, false).String() {
			return true, nil
		}
		if remote.String() == repo.URL(ctx, true).String() {
			return true, nil
		}
		return true, gogh.ErrProjectAlreadyExists
	}
	return false, nil
}
