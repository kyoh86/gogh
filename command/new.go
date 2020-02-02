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
	gitClient *GitClient,
	hubClient *HubClient,
	private bool,
	description string,
	homepage *url.URL,
	bare bool,
	template string,
	separateGitDir string,
	shared gogh.ProjectShared,
	repo *gogh.Repo,
) error {
	log.Printf("info: Creating new project and a remote repository %s", repo)
	project, err := gogh.FindOrNewProject(ctx, repo)
	if err != nil {
		return err
	}
	if project.Exists {
		return gogh.ErrProjectAlreadyExists
	}

	// mkdir
	log.Println("info: Creating a directory")
	if err := os.MkdirAll(project.FullPath, os.ModePerm); err != nil {
		return err
	}

	// git init
	log.Println("info: Initializing a repository")
	if err := gitClient.Init(ctx, project, bare, template, separateGitDir, shared); err != nil {
		return err
	}

	// hub create
	log.Println("info: Creating a new repository in GitHub")
	if _, err := hubClient.Create(ctx, repo, description, homepage, private); err != nil {
		return err
	}

	return nil
}
