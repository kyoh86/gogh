package gogh

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/repo"
)

// Pipe handles like `gogh pipe github-list-starred kyoh86` calling `github-list-starred kyoh86` and bulk its output
func Pipe(update, withSSH, shallow bool, command string, commandArgs []string) error {
	cmd := exec.Command(command, commandArgs...)
	cmd.Stderr = os.Stderr

	in, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	defer cmd.Wait()
	return bulkFromReader(in, update, withSSH, shallow)
}

// Bulk get repositories specified in stdin.
func Bulk(update, withSSH, shallow bool) error {
	return bulkFromReader(os.Stdin, update, withSSH, shallow)
}

// bulkFromReader bulk get repositories specified in reader.
func bulkFromReader(in io.Reader, update, withSSH, shallow bool) error {
	var repoSpecs repo.Specs
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if err := repoSpecs.Set(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return GetAll(update, withSSH, shallow, repoSpecs)
}
