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
	private bool,
	description string,
	homepage *url.URL,
	browse bool,
	clipboard bool,
	bare bool,
	template string,
	separateGitDir string,
	shared gogh.ProjectShared,
	repo *gogh.Repo,
) error {
	log.Printf("info: Creating new project and a remote repository %s", repo)
	project, err := gogh.FindProject(ctx, repo)
	switch err {
	case gogh.ProjectNotFound:
		project, err = gogh.NewProject(ctx, repo)
		if err != nil {
			return err
		}
	case nil:
		return gogh.ProjectAlreadyExists
	default:
		return err
	}

	// mkdir
	log.Println("info: Creating a directory")
	if err := os.MkdirAll(project.FullPath, os.ModePerm); err != nil {
		return err
	}

	// git init
	log.Println("info: Initializing a repository")
	if err := gitInit(ctx, bare, template, separateGitDir, shared, project.FullPath); err != nil {
		return err
	}

	// hub create
	log.Println("info: Creating a new repository in GitHub")
	return hubCreate(ctx, private, description, homepage, browse, clipboard, repo, project.FullPath)
}
