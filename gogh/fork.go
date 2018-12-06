package gogh

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/github/hub/commands"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx Context, update, withSSH, shallow, noRemote bool, remoteName string, organization string, repoSpec RepoSpec) error {
	log.Printf("info: cloning a repository")
	if err := Get(ctx, update, withSSH, shallow, repoSpec); err != nil {
		return err
	}

	repo, err := FromURL(ctx, repoSpec.URL(ctx, withSSH))
	if err != nil {
		return err
	}
	if err := os.Chdir(repo.FullPath); err != nil {
		return err
	}

	log.Printf("info: forking a repository")
	var hubArgs []string
	hubArgs = appendIf(hubArgs, "--no-remote", noRemote)
	hubArgs = appendIfFilled(hubArgs, "--remote-name", remoteName)
	hubArgs = appendIfFilled(hubArgs, "--organization", organization)
	// call hub fork
	//UNDONE: Should I set GITHUB_HOST and HUB_PROTOCOL? : see `man hub`.
	log.Printf("debug: calling `hub fork %s`", strings.Join(hubArgs, " "))
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("fork"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	if _, err := fmt.Fprintln(ctx.Stdout(), repo.RelPath); err != nil {
		return err
	}
	return nil
}
