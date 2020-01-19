package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func (c *GitClient) GetRemotes(ctx gogh.Context, project *gogh.Project) (map[string]*url.URL, error) {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Stdin = ctx.Stdin()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = project.FullPath

	if err := execCommand(cmd); err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}

	remotes := map[string]*url.URL{}
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		u, err := url.Parse(strings.TrimSpace(fields[1]))
		if err != nil {
			return nil, err
		}
		remotes[fields[0]] = u
	}
	return remotes, nil
}

func (c *GitClient) GetRemote(ctx gogh.Context, project *gogh.Project, name string) (*url.URL, error) {
	cmd := exec.Command("git", "remote", "get-url", name)
	cmd.Stdin = ctx.Stdin()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = project.FullPath

	if err := execCommand(cmd); err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}
	return url.Parse(strings.TrimSpace(stdout.String()))
}

func (c *GitClient) RemoveRemote(ctx gogh.Context, project *gogh.Project, name string) error {
	cmd := exec.Command("git", "remote", "remove", name)
	cmd.Stdin = ctx.Stdin()
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	cmd.Dir = project.FullPath

	return execCommand(cmd)
}

func (c *GitClient) RenameRemote(ctx gogh.Context, project *gogh.Project, oldName, newName string) error {
	cmd := exec.Command("git", "remote", "rename", oldName, newName)
	cmd.Stdin = ctx.Stdin()
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	cmd.Dir = project.FullPath

	return execCommand(cmd)
}

func (c *GitClient) AddRemote(ctx gogh.Context, project *gogh.Project, name string, url *url.URL) error {
	cmd := exec.Command("git", "remote", "add", name, url.String())
	cmd.Stdin = ctx.Stdin()
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	cmd.Dir = project.FullPath

	return execCommand(cmd)
}
