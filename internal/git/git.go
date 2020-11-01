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

	"github.com/kyoh86/gogh/internal/delegate"
)

type Client struct {
}

func (c *Client) command(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
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

func (c *Client) Fetch(local string) error {
	cmd := c.command("fetch", "--all")
	cmd.Dir = local

	return delegate.ExecCommand(cmd)
}

func (c *Client) GetCurrentBranch(local string) (string, error) {
	cmd := c.command("branch", "--show-current")
	cmd.Dir = local

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = local

	if err := delegate.ExecCommand(cmd); err != nil {
		return "", fmt.Errorf("%s: %w", stderr.String(), err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (c *Client) SetUpstreamTo(local string, upstream string) error {
	cmd := c.command("branch", "--set-upstream-to", upstream)
	cmd.Dir = local

	return delegate.ExecCommand(cmd)
}

func (c *Client) Status(path string, out, err io.Writer) error {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stdout = out
	cmd.Stderr = err
	cmd.Dir = path
	return delegate.ExecCommand(cmd)
}

func (c *Client) GetStatusSummary(path string, err io.Writer) (StatusSummary, error) {
	parser := &porcelainCollector{}
	if err := c.Status(path, parser, err); err != nil {
		return StatusSummaryClear, err
	}
	return parser.Close(), nil
}

type StatusSummary int

const (
	StatusSummaryClear     = StatusSummary(0)
	StatusSummaryModified  = StatusSummary(1)
	StatusSummaryUntracked = StatusSummary(2)
)

func (c StatusSummary) String() string {
	switch c {
	default: // case StatusSummaryClear:
		return "  "
	case StatusSummaryUntracked:
		return " +"
	case StatusSummaryModified:
		return "M "
	case StatusSummaryUntracked | StatusSummaryModified:
		return "M+"
	}
}

type porcelainCollector struct {
	summary StatusSummary
	surplus string
}

func runeToCode(r rune) StatusSummary {
	switch r {
	case ' ':
		return StatusSummaryClear
	case '?':
		return StatusSummaryUntracked
	default:
		return StatusSummaryModified
	}
}

func (w *porcelainCollector) parseLine(line string) {
	if len(line) > 2 {
		w.summary |= runeToCode([]rune(line)[0])
		w.summary |= runeToCode([]rune(line)[1])
	}
}

func (w *porcelainCollector) Write(p []byte) (int, error) {
	lines := strings.Split(w.surplus+string(p), "\n")
	last := len(lines) - 1
	for i := 0; i < last; i++ {
		w.parseLine(lines[i])
	}
	w.surplus = lines[last]
	return len(p), nil
}

func (w *porcelainCollector) Close() StatusSummary {
	w.parseLine(w.surplus)
	summary := w.summary
	w.surplus = ""
	w.summary = StatusSummaryClear
	return summary
}
