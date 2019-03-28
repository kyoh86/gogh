package internal

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kyoh86/gogh/gogh"
)

type GitClient struct{}

func (c *GitClient) Init(
	ctx gogh.Context,
	project *gogh.Project,
	bare bool,
	template string,
	separateGitDir string,
	shared gogh.ProjectShared,
) error {
	args := []string{"init"}
	args = appendIf(args, "--bare", bare)
	args = appendIfFilled(args, "--template", template)
	args = appendIfFilled(args, "--separate-git-dir", separateGitDir)
	args = appendIfFilled(args, "--shared", shared.String())
	args = append(args, project.FullPath)
	cmd := exec.Command("git", args...)
	cmd.Stdin = ctx.Stdin()
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	return execCommand(cmd)
}

func (c *GitClient) Clone(ctx gogh.Context, project *gogh.Project, remote *url.URL, shallow bool) error {
	dir, _ := filepath.Split(project.FullPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth", "1")
	}
	args = append(args, remote.String(), project.FullPath)

	cmd := exec.Command("git", args...)
	cmd.Stdin = ctx.Stdin()
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	return execCommand(cmd)
}

func (c *GitClient) Update(ctx gogh.Context, project *gogh.Project) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Stdin = ctx.Stdin()
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	cmd.Dir = project.FullPath

	return execCommand(cmd)
}
