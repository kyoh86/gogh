package gogh

import (
	"fmt"
	"os"

	"github.com/github/hub/commands"
	"github.com/kyoh86/gogh/repo"
)

// Fork clone/sync with a remote repository make a fork of a remote repository on GitHub and add GitHub as origin
func Fork(update, withSSH, shallow, noRemote bool, remoteName string, organization string, repoSpec repo.Spec) error {
	if err := Get(update, withSSH, shallow, repoSpec); err != nil {
		return err
	}

	rmt, err := repoSpec.Remote(withSSH)
	if err != nil {
		return err
	}
	remoteURL := rmt.URL()
	local, err := repo.FromURL(remoteURL)
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
