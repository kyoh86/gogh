package gogh

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

func gitInit(
	ctx Context,
	bare bool,
	template string,
	separateGitDir string,
	shared RepoShared,
	directory string,
) error {
	args := []string{"init"}
	args = appendIf(args, "--bare", bare)
	args = appendIfFilled(args, "--template", template)
	args = appendIfFilled(args, "--separate-git-dir", separateGitDir)
	args = appendIfFilled(args, "--shared", shared.String())
	args = append(args, directory)
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	err := execCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

// gitClone git repository
func gitClone(ctx Context, remote *url.URL, local string, shallow bool) error {
	dir, _ := filepath.Split(local)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth", "1")
	}
	args = append(args, remote.String(), local)

	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	return execCommand(cmd)
}

// gitUpdate pulls changes from remote repository
func gitUpdate(ctx Context, local string) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Stdin = os.Stdin
	cmd.Stdout = ctx.Stdout()
	cmd.Stderr = ctx.Stderr()
	cmd.Dir = local

	return execCommand(cmd)
}
