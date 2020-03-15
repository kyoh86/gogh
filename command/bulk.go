package command

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/gogh"
)

// Pipe handles like `gogh pipe github-list-starred kyoh86` calling `github-list-starred kyoh86` and bulk its output
func Pipe(ev gogh.Env, gitClient GitClient, update, withSSH, shallow bool, command string, commandArgs []string) (retErr error) {
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
	return bulkFromReader(ev, gitClient, in, update, withSSH, shallow)
}

// Bulk get repositories specified in stdin.
func Bulk(ev gogh.Env, gitClient GitClient, update, withSSH, shallow bool) error {
	return bulkFromReader(ev, gitClient, os.Stdin, update, withSSH, shallow)
}

// bulkFromReader bulk get repositories specified in reader.
func bulkFromReader(ev gogh.Env, gitClient GitClient, in io.Reader, update, withSSH, shallow bool) error {
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

	repos, err := specs.Validate(ev)
	if err != nil {
		return err
	}
	return GetAll(ev, gitClient, update, withSSH, shallow, repos)
}
