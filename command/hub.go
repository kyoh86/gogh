package command

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/github/hub/commands"
	"github.com/kyoh86/gogh/gogh"
)

func hubFork(
	ctx gogh.Context,
	project *gogh.Project,
	remote *gogh.Remote,
	noRemote bool,
	remoteName string,
	organization string,
) (retErr error) {
	cd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(project.FullPath); err != nil {
		return err
	}
	defer func() {
		if err := os.Chdir(cd); err != nil && retErr == nil {
			retErr = err
		}
	}()

	var hubArgs []string
	hubArgs = appendIf(hubArgs, "--no-remote", noRemote)
	hubArgs = appendIfFilled(hubArgs, "--remote-name", remoteName)
	hubArgs = appendIfFilled(hubArgs, "--organization", organization)
	// call hub fork
	if err := os.Setenv("GITHUB_HOST", remote.Host(ctx)); err != nil {
		return err
	}
	log.Printf("debug: calling `hub fork %s`", strings.Join(hubArgs, " "))
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
	remote *gogh.Remote,
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
	hubArgs = appendIf(hubArgs, "-p", private)
	hubArgs = appendIf(hubArgs, "-o", browse)
	hubArgs = appendIf(hubArgs, "-c", clipboard)
	hubArgs = appendIfFilled(hubArgs, "-d", description)
	if homepage != nil {
		hubArgs = append(hubArgs, "-h", homepage.String())
	}
	hubArgs = append(hubArgs, remote.URL(ctx, false).String())
	log.Printf("debug: calling `hub create %s`", strings.Join(hubArgs, " "))
	if err := os.Setenv("GITHUB_HOST", remote.Host(ctx)); err != nil {
		return err
	}
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("create"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	return nil
}
