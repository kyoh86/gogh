package command

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/github/hub/commands"
	"github.com/kyoh86/gogh/gogh"
)

func chdirTmp(dir string) (func() error, error) {
	log.Printf("debug: Changing working directory to %s", dir)
	cd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(dir); err != nil {
		return nil, err
	}
	return func() error {
		log.Printf("debug: Changing back working directory to %s", cd)
		return os.Chdir(cd)
	}, nil
}

func hubFork(
	ctx gogh.Context,
	project *gogh.Project,
	repo *gogh.Repo,
	noRemote bool,
	remoteName string,
	organization string,
) (retErr error) {
	tear, err := chdirTmp(project.FullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := tear(); err != nil && retErr == nil {
			retErr = err
		}
	}()

	var hubArgs []string
	hubArgs = appendIf(hubArgs, "--no-remote", noRemote)
	hubArgs = appendIfFilled(hubArgs, "--remote-name", remoteName)
	hubArgs = appendIfFilled(hubArgs, "--organization", organization)
	// call hub fork
	if err := os.Setenv("GITHUB_HOST", repo.Host(ctx)); err != nil {
		return err
	}
	if err := os.Setenv("GITHUB_TOKEN", ctx.GitHubToken()); err != nil {
		return err
	}
	log.Printf("debug: Calling `hub fork %s`", strings.Join(hubArgs, " "))
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("fork"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}

	return nil
}

func hubCreate(
	ctx gogh.Context,
	private bool,
	description string,
	homepage *url.URL,
	browse bool,
	clipboard bool,
	repo *gogh.Repo,
	directory string,
) (retErr error) {
	tear, err := chdirTmp(directory)
	if err != nil {
		return err
	}
	defer func() {
		if err := tear(); err != nil && retErr == nil {
			retErr = err
		}
	}()

	var hubArgs []string
	hubArgs = appendIf(hubArgs, "-p", private)
	hubArgs = appendIf(hubArgs, "-o", browse)
	hubArgs = appendIf(hubArgs, "-c", clipboard)
	hubArgs = appendIfFilled(hubArgs, "-d", description)
	if homepage != nil {
		hubArgs = append(hubArgs, "-h", homepage.String())
	}
	hubArgs = append(hubArgs, repo.URL(ctx, false).String())
	log.Printf("debug: Calling `hub create %s`", strings.Join(hubArgs, " "))
	if err := os.Setenv("GITHUB_HOST", repo.Host(ctx)); err != nil {
		return err
	}
	if err := os.Setenv("GITHUB_TOKEN", ctx.GitHubToken()); err != nil {
		return err
	}
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("create"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	return nil
}
