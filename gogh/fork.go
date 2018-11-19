package gogh

import (
	"fmt"
	"os"

	"github.com/github/hub/commands"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(ctx Context, update, withSSH, shallow, noRemote bool, remoteName string, organization string, repoSpec RepoSpec) error {
	if err := Get(ctx, update, withSSH, shallow, repoSpec); err != nil {
		return err
	}

	rmt, err := repoSpec.Remote(ctx, withSSH)
	if err != nil {
		return err
	}
	remoteURL := rmt.URL()
	local, err := FromURL(ctx, remoteURL)
	if err != nil {
		return err
	}
	if err := os.Chdir(local.FullPath); err != nil {
		return err
	}
	var hubArgs []string
	hubArgs = appendIf(hubArgs, "--no-remote", noRemote)
	hubArgs = appendIfFilled(hubArgs, "--remote-name", remoteName)
	hubArgs = appendIfFilled(hubArgs, "--organization", organization)
	// call hub fork
	//UNDONE: Should I set GITHUB_HOST and HUB_PROTOCOL? : see `man hub`.
	execErr := commands.CmdRunner.Call(commands.CmdRunner.Lookup("fork"), commands.NewArgs(hubArgs))
	if execErr.Err != nil {
		return execErr.Err
	}
	fmt.Println(local.RelPath)
	return nil
}
