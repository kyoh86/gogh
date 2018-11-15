package gogh

import (
	"net/url"
	"os"

	"github.com/github/hub/commands"
	"github.com/kyoh86/gogh/internal/run"
	"github.com/kyoh86/gogh/repo"
)

func hubInit(
	bare bool,
	template string,
	separateGitDir string,
	shared repo.Shared,
	directory string,
) error {
	var hubArgs []string
	if bare {
		hubArgs = append(hubArgs, "--bare")
	}
	if template != "" {
		hubArgs = append(hubArgs, "--template", template)
	}
	if separateGitDir != "" {
		hubArgs = append(hubArgs, "--separate-git-dir", separateGitDir)
	}
	if shared != "" {
		hubArgs = append(hubArgs, "--shared", shared.String())
	}
	hubArgs = append(hubArgs, directory)
	//UNDONE: Should I set GITHUB_HOST and HUB_PROTOCOL? : see `man hub`.
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("init"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	return nil
}

func hubCreate(
	private bool,
	description string,
	homepage *url.URL,
	browse bool,
	clipboard bool,
	repoName repo.Name,
	directory string,
) (retErr error) {
	// cd
	cd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(directory); err != nil {
		return err
	}
	defer func() {
		if err := os.Chdir(cd); err != nil && retErr == nil {
			retErr = err
		}
	}()

	var hubArgs []string
	if private {
		hubArgs = append(hubArgs, "-p")
	}
	if description != "" {
		hubArgs = append(hubArgs, "-d", description)
	}
	if browse {
		hubArgs = append(hubArgs, "-o")
	}
	if clipboard {
		hubArgs = append(hubArgs, "-c")
	}
	hubArgs = append(hubArgs, repoName.String())
	//UNDONE: Should I set GITHUB_HOST and HUB_PROTOCOL? : see `man hub`.
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("create"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	return nil
}

func nameToPath(name string) (string, error) {
	spec, err := repo.NewSpec(name)
	if err != nil {
		return "", err
	}
	rem, err := spec.Remote(false)
	if err != nil {
		return "", err
	}
	loc, err := repo.FromURL(rem.URL())
	if err != nil {
		return "", err
	}
	return loc.FullPath, nil
}

// New creates a repository in local and remote.
func New(
	private bool,
	description string,
	homepage *url.URL,
	browse bool,
	clipboard bool,
	bare bool,
	template string,
	separateGitDir string,
	shared repo.Shared,
	repoName repo.Name,
) error {
	path, err := nameToPath(repoName.String())
	if err != nil {
		return err
	}

	// mkdir
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	// hub init
	if err := hubInit(bare, template, separateGitDir, shared, path); err != nil {
		return err
	}

	// hub create
	if err := hubCreate(private, description, homepage, browse, clipboard, repoName, path); err != nil {
		return err
	}

	// which yo
	if err := run.Silently("which", "yo"); err == nil {
		if err := run.InDir("yo", path); err != nil {
			return err
		}
	}
	return nil
}
