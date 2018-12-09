package gogh

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/github/hub/commands"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx Context, update, withSSH, shallow, noRemote bool, remoteName string, organization string, remote *Remote) error {
	log.Printf("info: cloning a repository")
	if err := Get(ctx, update, withSSH, shallow, remote); err != nil {
		return err
	}

	local, err := FindLocal(ctx, remote)
	if err != nil {
		return err
	}
	if err := os.Chdir(local.FullPath); err != nil {
		return err
	}

	log.Printf("info: forking a repository")
	var hubArgs []string
	hubArgs = appendIf(hubArgs, "--no-remote", noRemote)
	hubArgs = appendIfFilled(hubArgs, "--remote-name", remoteName)
	hubArgs = appendIfFilled(hubArgs, "--organization", organization)
	// call hub fork
	os.Setenv("GITHUB_HOST", remote.Host(ctx))
	log.Printf("debug: calling `hub fork %s`", strings.Join(hubArgs, " "))
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("fork"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	if _, err := fmt.Fprintln(ctx.Stdout(), local.RelPath); err != nil {
		return err
	}
	return nil
}
