package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gogh/internal/delegate"
)

func New(ctx gogh.IOContext) *Client {
	return &Client{
		For: struct {
			Stdout io.Writer
			Stderr io.Writer
			Stdin  io.Reader
		}{
			ctx.Stdout(),
			ctx.Stderr(),
			ctx.Stdin(),
		},
	}
}

type Client struct {
	For struct {
		Stdout io.Writer
		Stderr io.Writer
		Stdin  io.Reader
	}
}

func (c *Client) command(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Stdin = c.For.Stdin
	cmd.Stdout = c.For.Stdout
	cmd.Stderr = c.For.Stderr
	return cmd
}

func (c *Client) Init(
	directory string,
	bare bool,
	template string,
	separateGitDir string,
	shared string,
) error {
	args := []string{"init"}
	args = delegate.AppendIf(args, "--bare", bare)
	args = delegate.AppendIfFilled(args, "--template", template)
	args = delegate.AppendIfFilled(args, "--separate-git-dir", separateGitDir)
	args = delegate.AppendPairedIfFilled(args, "--shared", shared)
	args = append(args, directory)
	return delegate.ExecCommand(c.command(args...))
}

func (c *Client) Clone(local string, remote *url.URL, shallow bool) error {
	dir, _ := filepath.Split(local)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth", "1")
	}
	args = append(args, remote.String(), local)

	cmd := c.command(args...)
	return delegate.ExecCommand(cmd)
}

func (c *Client) Update(local string) error {
	cmd := c.command("pull", "--ff-only")
	cmd.Dir = local

	return delegate.ExecCommand(cmd)
}

func (c *Client) GetRemotes(local string) (map[string]*url.URL, error) {
	cmd := c.command("remote", "-v")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = local

	if err := delegate.ExecCommand(cmd); err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}

	remotes := map[string]*url.URL{}
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		u, err := url.Parse(strings.TrimSpace(fields[1]))
		if err != nil {
			return nil, err
		}
		remotes[fields[0]] = u
	}
	return remotes, nil
}

func (c *Client) GetRemote(local string, name string) (*url.URL, error) {
	cmd := c.command("remote", "get-url", name)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = local

	if err := delegate.ExecCommand(cmd); err != nil {
		return nil, fmt.Errorf("%s: %w", stderr.String(), err)
	}
	return url.Parse(strings.TrimSpace(stdout.String()))
}

func (c *Client) RemoveRemote(local string, name string) error {
	cmd := c.command("remote", "remove", name)
	cmd.Dir = local

	return delegate.ExecCommand(cmd)
}

func (c *Client) RenameRemote(local string, oldName, newName string) error {
	cmd := c.command("remote", "rename", oldName, newName)
	cmd.Dir = local

	return delegate.ExecCommand(cmd)
}

func (c *Client) AddRemote(local string, name string, url *url.URL) error {
	cmd := c.command("remote", "add", name, url.String())
	cmd.Dir = local

	return delegate.ExecCommand(cmd)
}
