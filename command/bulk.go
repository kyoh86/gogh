package command

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/gogh"
)

// Pipe handles like `gogh pipe github-list-starred kyoh86` calling `github-list-starred kyoh86` and bulk its output
func Pipe(env gogh.Env, gitClient GitClient, update, withSSH, shallow bool, command string, commandArgs []string) (retErr error) {
	cmd := exec.Command(command, commandArgs...)
	in, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	defer func() {
		if err := cmd.Wait(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return bulkFromReader(env, gitClient, in, update, withSSH, shallow)
}

// Bulk get repositories specified in stdin.
func Bulk(env gogh.Env, gitClient GitClient, update, withSSH, shallow bool) error {
	return bulkFromReader(env, gitClient, os.Stdin, update, withSSH, shallow)
}

// bulkFromReader bulk get repositories specified in reader.
func bulkFromReader(env gogh.Env, gitClient GitClient, in io.Reader, update, withSSH, shallow bool) error {
	var specs gogh.RepoSpecs
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if err := specs.Set(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	repos, err := specs.Validate(env)
	if err != nil {
		return err
	}
	return GetAll(env, gitClient, update, withSSH, shallow, repos)
}
