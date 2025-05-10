package view

import (
	"encoding/json"
	"strings"

	"github.com/kyoh86/gogh/v3/core/workspace"
)

// LocalRepoFormat defines the interface for formatting local repository references
type LocalRepoFormat interface {
	Format(ref workspace.RepoInfo) (string, error)
}

// LocalRepoFormatFunc is a function type that implements the LocalRepoFormat interface
type LocalRepoFormatFunc func(workspace.RepoInfo) (string, error)

// Format calls the function itself to format the local repository reference
func (f LocalRepoFormatFunc) Format(ref workspace.RepoInfo) (string, error) {
	return f(ref)
}

// LocalRepoFormatRelPath formats the local repository reference to its full path
var LocalRepoFormatFullPath = LocalRepoFormatFunc(func(ref workspace.RepoInfo) (string, error) {
	return ref.FullPath(), nil
})

// LocalRepoFormatRelFilePath formats the local repository reference to its path
var LocalRepoFormatPath = LocalRepoFormatFunc(func(ref workspace.RepoInfo) (string, error) {
	return ref.Path(), nil
})

// LocalRepoFormatJSON formats the local repository reference to a JSON string
var LocalRepoFormatJSON = LocalRepoFormatFunc(func(ref workspace.RepoInfo) (string, error) {
	buf, _ := json.Marshal(map[string]any{
		"fullPath": ref.FullPath(),
		"path":     ref.Path(),
		"host":     ref.Host(),
		"owner":    ref.Owner(),
		"name":     ref.Name(),
	})
	return string(buf), nil
})

// LocalRepoFormatFields formats the local repository reference to a string with specified fields
func LocalRepoFormatFields(s string) LocalRepoFormat {
	return LocalRepoFormatFunc(func(ref workspace.RepoInfo) (string, error) {
		return strings.Join([]string{
			ref.FullPath(),
			ref.Path(),
			ref.Host(),
			ref.Owner(),
			ref.Name(),
		}, s), nil
	})
}
