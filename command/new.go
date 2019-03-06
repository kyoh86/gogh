package command

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"

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
	project, err := gogh.FindProject(ctx, repo)
	if err != nil {
		return err
	}

	// mkdir
	log.Println("info: creating a directory")
	if err := os.MkdirAll(project.FullPath, os.ModePerm); err != nil {
		return err
	}

	// git init
	log.Println("info: initializing a repository")
	if err := gitInit(ctx, bare, template, separateGitDir, shared, project.FullPath); err != nil {
		return err
	}

	// hub create
	log.Println("info: creating a new repository in GitHub")
	if err := hubCreate(ctx, private, description, homepage, browse, clipboard, repo, project.FullPath); err != nil {
		return err
	}

	// which yo
	cmd := exec.Command("which", "yo")
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	if err := execCommand(cmd); err == nil {
		log.Println("info: calling yo")
		cmd := exec.Command("yo")
		cmd.Stdin = os.Stdin
		cmd.Stdout = ctx.Stdout()
		cmd.Stderr = ctx.Stderr()
		cmd.Dir = project.FullPath
		if err := execCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}
