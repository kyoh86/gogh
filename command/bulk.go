package command

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/kyoh86/gogh/gogh"
)

// Pipe handles like `gogh pipe github-list-starred kyoh86` calling `github-list-starred kyoh86` and bulk its output
func Pipe(ctx gogh.Context, gitClient GitClient, update, withSSH, shallow bool, command string, commandArgs []string) (retErr error) {
	InitLog(ctx)

	cmd := exec.Command(command, commandArgs...)
	if ctx, ok := ctx.(gogh.IOContext); ok {
		cmd.Stderr = ctx.Stderr()
	}

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
	return bulkFromReader(ctx, gitClient, in, update, withSSH, shallow)
}

// Bulk get repositories specified in stdin.
func Bulk(ctx gogh.Context, gitClient GitClient, update, withSSH, shallow bool) error {
	InitLog(ctx)
	var stdin io.Reader = os.Stdin
	if ctx, ok := ctx.(gogh.IOContext); ok {
		stdin = ctx.Stdin()
	}
	return bulkFromReader(ctx, gitClient, stdin, update, withSSH, shallow)
}

// bulkFromReader bulk get repositories specified in reader.
func bulkFromReader(ctx gogh.Context, gitClient GitClient, in io.Reader, update, withSSH, shallow bool) error {
	var repos gogh.Repos
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if err := repos.Set(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return GetAll(ctx, gitClient, update, withSSH, shallow, repos)
}
