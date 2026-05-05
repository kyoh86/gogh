package worktree

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format defines the interface for formatting worktrees to a string representation
type Format interface {
	Format(wt Worktree, repo string) (string, error)
}

// FormatFunc is a function type that implements the Format interface
type FormatFunc func(wt Worktree, repo string) (string, error)

// Format calls the function itself to format the worktree
func (f FormatFunc) Format(wt Worktree, repo string) (string, error) {
	return f(wt, repo)
}

// FormatDefault is the default format: repository name with branch in parentheses
var FormatDefault = FormatFunc(func(wt Worktree, repo string) (string, error) {
	if wt.Branch == "" {
		return fmt.Sprintf("%s (detached)", repo), nil
	}
	return fmt.Sprintf("%s (%s)", repo, wt.Branch), nil
})

// FormatFullPath formats the worktree to its full path
var FormatFullPath = FormatFunc(func(wt Worktree, repo string) (string, error) {
	return wt.Path, nil
})

// FormatJSON formats the worktree to a JSON string
var FormatJSON = FormatFunc(func(wt Worktree, repo string) (string, error) {
	buf, err := json.Marshal(map[string]any{
		"repo":   repo,
		"path":   wt.Path,
		"branch": wt.Branch,
		"commit": wt.Commit,
	})
	if err != nil {
		return "", err
	}
	return string(buf), nil
})

// FormatFields formats the worktree with specified fields
func FormatFields(separator string) Format {
	return FormatFunc(func(wt Worktree, repo string) (string, error) {
		return strings.Join([]string{
			wt.Path,
			repo,
			wt.Branch,
			wt.Commit,
		}, separator), nil
	})
}
