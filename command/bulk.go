package command

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/gogh"
)

// Pipe handles like `gogh pipe github-list-starred kyoh86` calling `github-list-starred kyoh86` and bulk its output
func Pipe(ctx gogh.Context, update, withSSH, shallow bool, command string, commandArgs []string) (retErr error) {
	cmd := exec.Command(command, commandArgs...)
	cmd.Stderr = ctx.Stderr()

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
	return bulkFromReader(ctx, in, update, withSSH, shallow)
}

// Bulk get repositories specified in stdin.
func Bulk(ctx gogh.Context, update, withSSH, shallow bool) error {
	return bulkFromReader(ctx, os.Stdin, update, withSSH, shallow)
}

// bulkFromReader bulk get repositories specified in reader.
func bulkFromReader(ctx gogh.Context, in io.Reader, update, withSSH, shallow bool) error {
	var remotes gogh.Remotes
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if err := remotes.Set(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return GetAll(ctx, update, withSSH, shallow, remotes)
}